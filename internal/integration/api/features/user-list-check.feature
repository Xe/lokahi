Feature: A user can list Checks

  Background:
    Given a base stack
    Given I want to list checks
    # make many checks
    Given an example check
    Given an example check
    Given an example check
    Given an example check
    Given an example check
    Given an example check
    Given an example check
    Given an example check
    Given an example check
    Given I can fetch the check

  Scenario: list checks
    As a User
    When I try to list checks
    Then there was no error
    Then tear everything down

  Scenario: guaranteed page 2
    As a User
    Given check list count is 1
    Given check list offset is 1
    When I try to list checks
    Then there was no error
    Then tear everything down
