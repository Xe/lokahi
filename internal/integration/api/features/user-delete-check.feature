Feature: A User can delete a Check

  Background:
    Given a base stack
    Given an example check
    Given I can fetch the check

  Scenario: normal check deletion
    As a User
    When I try to delete the check
    Then there was no error
    Then tear everything down

  Scenario: cannot delete a check that does not exist
    As a User
    Given a random check ID
    When I try to delete the check
    Then there was an error
    Then tear everything down
