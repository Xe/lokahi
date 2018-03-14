Feature: A User can delete a Check

  Background:
    Given a base stack
    Given an example check
    Given I can fetch the check

  Scenario: normal check deletion
    As a User
    When I delete the check
    Then the check cannot be fetched

  Scenario: cannot delete a check that does not exist
    As a User
    Given a random check ID
    When I cant delete the check
    Then the check cannot be fetched
