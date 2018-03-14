Feature: A User can create a Check

  Scenario: normal check creation
    As a User
    Given a base stack
    Given I want to create a check
    Given a check monitoring url of "https://google.com"
    Given a check webhook url of "http://sample_hook:9001/twirp/github.xe.lokahi.Webhook/Handle"
    Given a check every of 60
    Given a check playbook url of "https://figureit.out"
    When I try to create the check
    Then there was no error
    Then I can fetch the check
    Then tear everything down
