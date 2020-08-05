package main

// EmailMessage sends email by using inputs via ArgumentOptions
type EmailMessage struct {
	args ArgumentOptions
}

// Send message via email
func (email *EmailMessage) Send(args ArgumentOptions)  {
	email.args = args;
}