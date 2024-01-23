package mail

import (
	"fmt"
	"github.com/kanishkmittal55/simplebank/db/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSendEmailWithGmail(t *testing.T) {
	//if testing.Short() {
	//	t.Skip()
	//}
	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	// Print the config values to verify them
	//fmt.Printf("Email Sender Name: %s\n", config.EmailSenderName)
	//fmt.Printf("Email Sender Address: %s\n", config.EmailSenderAddress)
	//fmt.Printf("Email Sender Password: %s\n", config.EmailSenderPassword)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "A test Email"
	Content := `
	<h1> Test Email </h1>
	<p> This is a test message <a href="http://hassleskip.com">Click to Visit Website</a></p>
	`

	to := []string{
		"kanishkmittal55@gmail.com",
		"2156564@brunel.ac.uk",
	}

	fmt.Printf("Email Sender Address: %s", to)
	// attachFiles := []string{"../README.md"}
	err = sender.SendEmail(subject, Content, to, nil, nil, nil)
	require.NoError(t, err)
}
