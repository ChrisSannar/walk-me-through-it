# Notes

Here's the place to put the human readable ideas, potential features, and goals for the project.

## Next

  - Include a way to delete walkthroughs
  - Review security flaws in code structure as well as potential vulnerabilities
  - Consider how the code can be refactored. Ask the AI to build a walkthrough for what we could do, and why we would do it.
  - Fix the screen flickering bug
  - Delicately add more tests before the project gets too big
  - Layered walkthroughs?
    * You're walking through one thing, but then you want to learn about something else, so you create a new walkthrough and add it to the "stack". Once you're done with that walk through, you give it points and continue on.

## Goals

### MVP

 - Allow users to write and select their own `*.wmti.json` files and run them in the project (write and provide a style guide so that an LLM can generate them. Is that a "Skill"?)
 - Open and navigate through entire projects.

## TODO

 - Research how editors scan the structure of a file and how to map that to an AI.
 - How do I get an AI to do a better analysis? Do I have multiple AI's generate summaries based on what they know? How do I get them to create a decent context that they can work?

## Extra Features

 - A chat box inside the "Content" view that the user can ask about things in the project.
 - Have the option to "Scan" for new `*.wmti.json` files
 - When starting with the init, auto-detect AI keys in `.bashrc` (or congruent operating system) and then ask to use them

## Engineering

 - Maybe keep all the `*.wmti.json` files somewhere in the local storage file so they don't clutter things up. (also include an option to store them locally if the project is to be shared)
