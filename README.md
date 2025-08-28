# Typer ✍️

Typer is a lightweight typing practice program written in Go.  
It combines a simple CLI game, a quotes JSON store, and a playful JSON API.  
Think of it as a minimal offline **typeracer** – no opponents, no scoreboards, just you and the keyboard.

---

## Features

- **CLI Typing Game**  
  - Presents you with a random quote to type out.  
  - Tracks your accuracy and speed (basic metrics, extensible later).  

- **Quotes Loader**  
  - `quote.go` loads quotes from a JSON file.  
  - Easily extendable with your own quotes.  

- **JSON API**  
  - Serves quotes via a lightweight API.  
  - Useful for experimenting with Go’s net/http, or just for fun.

---

## Installation

Clone the repository:

```bash
git clone https://github.com/yourusername/typer.git
cd typer
