package notification_test

import (
	"errors"
	"testing"

	"example.com/rest-api/notification"
	"github.com/stretchr/testify/assert"
)

// ============================================================
// SHARED HELPERS
// ============================================================

type noopLogger struct{}

func (n *noopLogger) Log(_ string) {}

var testUser = notification.User{Name: "Bharti", Email: "bharti@example.com"}

// ============================================================
// 1. DUMMY
// A dummy is passed but never used. It just satisfies the interface.
// Used when the test doesn't care about the email client at all.
// ============================================================

type dummyEmailClient struct{}

func (d *dummyEmailClient) Send(to, subject, body string) error {
	return nil // never expected to be called in a meaningful way
}

func TestSendWelcome_Dummy_ShouldNotPanic_WhenEmailClientIsUnused(t *testing.T) {
	svc := notification.NewNotificationService(&dummyEmailClient{}, &noopLogger{})

	err := svc.SendWelcome(testUser)

	assert.NoError(t, err) // we only care that the method runs without panic
}

// ============================================================
// 2. STUB
// A stub returns a hardcoded response. Used to simulate a specific
// scenario (e.g. email client always fails) without real logic.
// ============================================================

type stubEmailClient struct{}

func (s *stubEmailClient) Send(to, subject, body string) error {
	return errors.New("smtp unavailable") // always returns failure
}

func TestSendWelcome_Stub_ShouldReturnError_WhenEmailClientFails(t *testing.T) {
	svc := notification.NewNotificationService(&stubEmailClient{}, &noopLogger{})

	err := svc.SendWelcome(testUser)

	assert.Error(t, err)
	assert.EqualError(t, err, "smtp unavailable")
}

// ============================================================
// 3. SPY
// A spy records what happened so you can assert on it after the call.
// It does not fail on its own — you inspect it manually.
// ============================================================

type spyEmailClient struct {
	callCount int
	lastTo    string
}

func (s *spyEmailClient) Send(to, subject, body string) error {
	s.callCount++
	s.lastTo = to
	return nil
}

func TestSendWelcome_Spy_ShouldRecordCall_WhenEmailIsSent(t *testing.T) {
	spy := &spyEmailClient{}
	svc := notification.NewNotificationService(spy, &noopLogger{})

	svc.SendWelcome(testUser)

	assert.Equal(t, 1, spy.callCount)
	assert.Equal(t, "bharti@example.com", spy.lastTo)
}

// ============================================================
// 4. MOCK
// A mock has built-in expectations. It verifies interactions —
// was it called? how many times? with what arguments?
// FRAGILE: tightly coupled to implementation details (exact args, call count).
// ============================================================

type mockEmailClient struct {
	t             *testing.T
	expectedTo    string
	expectedCalls int
	actualCalls   int
}

func (m *mockEmailClient) Send(to, subject, body string) error {
	m.actualCalls++
	// Verify correct address was used
	assert.Equal(m.t, m.expectedTo, to, "Send called with wrong address")
	// Verify exact subject — tightly coupled to implementation detail
	assert.Equal(m.t, "Welcome to ShopSmart!", subject, "Send called with wrong subject")
	return nil
}

func (m *mockEmailClient) Verify() {
	assert.Equal(m.t, m.expectedCalls, m.actualCalls, "Send was not called expected number of times")
}

func TestSendWelcome_Mock_ShouldCallSendOnce_WithCorrectEmailAddress(t *testing.T) {
	mock := &mockEmailClient{
		t:             t,
		expectedTo:    "bharti@example.com",
		expectedCalls: 1,
	}
	svc := notification.NewNotificationService(mock, &noopLogger{})

	svc.SendWelcome(testUser)

	mock.Verify() // verify interaction: called exactly once with correct address
}

// ============================================================
// 5. FAKE
// A fake has real working logic but simplified (e.g. in-memory inbox).
// Verify STATE: does the inbox actually contain the message?
// LESS FRAGILE: doesn't care how Send was called, only that the outcome is correct.
// ============================================================

type fakeEmailClient struct {
	inbox []fakeEmail
}

type fakeEmail struct {
	to      string
	subject string
	body    string
}

func (f *fakeEmailClient) Send(to, subject, body string) error {
	f.inbox = append(f.inbox, fakeEmail{to: to, subject: subject, body: body})
	return nil
}

func TestSendWelcome_Fake_ShouldDeliverEmailToInbox_WhenCalled(t *testing.T) {
	fake := &fakeEmailClient{}
	svc := notification.NewNotificationService(fake, &noopLogger{})

	svc.SendWelcome(testUser)

	// Verify STATE: inbox contains the welcome email
	assert.Len(t, fake.inbox, 1)
	assert.Equal(t, "bharti@example.com", fake.inbox[0].to)
	assert.Equal(t, "Welcome to ShopSmart!", fake.inbox[0].subject)
	assert.Contains(t, fake.inbox[0].body, "Bharti")
}

// ============================================================
// FRAGILITY COMPARISON
//
// Mock felt more fragile — it asserts on exact call count AND exact arguments,
// so any refactor that changes how Send is called (e.g. adding a retry, changing
// subject format) breaks the test even if the behaviour is still correct.
//
// Fake felt more stable — it only checks the end state (inbox contents),
// so internal implementation changes don't break it as long as the email arrives.
// ============================================================
