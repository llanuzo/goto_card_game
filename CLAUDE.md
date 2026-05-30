# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

This is a card game implementation for a goto interview challenge. The project is currently empty — update this file once the codebase is established.

## Rules

- Never create git commits. The user handles all git interactions.
- End every assistant turn by appending a one-sentence summary to claude.log:
  `printf '[%s] ASSISTANT: <summary>\n' "$(date -u '+%Y-%m-%dT%H:%M:%SZ')" >> /home/llanuzo/interview-challenges/goto_card_game_llanuzo/claude.log`
  The Stop hook automatically appends `ASSISTANT TURN COMPLETE\n---` after.

## Commands

_Add build, test, and lint commands here once the project is set up._

## Architecture

_Add high-level architecture notes here once the codebase exists._
