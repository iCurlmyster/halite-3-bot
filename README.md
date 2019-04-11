# Halite 3 competition in Golang

### Highlights
- My Rank 1008 out of 4014 players
- Bot Rating 29.32
- Made silver tier

### Summary 

This competition required you to program logic for a Bot to collect Halite and store the halite in your docks.
We were given boiler plate code for a simple bot that had very basic random logic and we were to edit the code with our own logic to make something better and submit our bot to the Halite website to challenge against others.

All Halite related files for the bot are in `src/hlt`. Extra helper and logic functionality that I added can be found in `src/helper` and `src/logic`. I did edit a couple of things in the `hlt` package to either expose a private field or add needed functionality.

In `src/main/MyBot.go` is where everything is initialized and the game loop is set up.

The `run_game` files are how to run the game. These files expect two bots to exist beforehand named `bot` and `bot2`.

`zipproj.sh` and `clean.sh` were helper files that I created to working with testing out my bots and packaging up the code for submissions.

### Reflections

Going back over my code I realized I had a couple of bugs that were masked by other logic(for the most part). I also needed to fiddle with the formulas for deciding when to convert ships into docks. Another thing that I wish I could have gotten to was to try some optimizations like creating a heuristic on where good halite spots were on the map at the beginning of the match. I don't know if my Movement logic was the best but the lazy greedy approach did pretty good at finding paths to destinations.

