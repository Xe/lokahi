Feature: A User can put updates into a Check

  Background:
    Given a base stack
    Given an example check

  Scenario: change url
    Given a random url in the last check
    When I try to put the check
    Then there was no error
    Then tear everything down
