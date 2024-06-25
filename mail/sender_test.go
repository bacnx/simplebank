package mail

import (
	"testing"

	"github.com/bacnx/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestSendMailWithGmail(t *testing.T) {
	config, err := util.GetConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.FromEmailName, config.FromEmailAddress, config.FromEmailPassword)

	content := `
    <h1>Simple Bank</h1>
    <p>This just a testing</p>
  `
	err = sender.SendEmail([]string{"nxbac.testing@gmail.com"}, "Test send email", content, nil, nil, []string{"../README.md"})
	require.NoError(t, err)
}
