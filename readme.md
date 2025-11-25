-----

# ssho 🚀

**The "Infinite Canvas" Terminal Multiplexer.**
*(Because `alt-tab` is for quitters.)*


Have you ever looked at your terminal and thought, *"This needs more... void"*? Or maybe you wanted to drag your `htop` window around like a sticky note on a whiteboard?

**ssho** is a local, web-based terminal emulator that treats your shell sessions like objects on an infinite canvas. It's like if **Tmux** went on a date with **Figma** and they had a beautiful, GPU-accelerated baby.

![image.jpg](static%2Fimage.jpg)

## ✨ Features (The Cool Stuff)

* **♾️ Infinite Canvas:** Pan, scroll, and explore the void. Your terminal windows live in a boundless world now.
* **🖱️ Drag & Drop:** Move windows around. Resize them. Organize your chaos.
* **🧠 Smart Repulsion:** New windows automatically push old ones out of the way. It's like they have personal space issues.
* **🧟 Zombie-Proof (Persistence):** Accidentally refreshed the page? Browser crashed? **Don't panic.** Your sessions (and your running `vim`) are still there waiting for you.
* **🎨 Eye-Candy Themes:** Comes with 5 hand-picked themes (GitHub Dark, One Dark, Solarized, Monokai, etc.). Your eyes will thank you.
* **🚀 GPU Accelerated:** Powered by `xterm-addon-webgl`. Typing feels smoother than butter on a hot GPU.
* **📦 Single Binary:** No dependencies. No `node_modules` black hole. Just one binary file to rule them all.

## 🛠️ Installation

### The "I Trust No One" Way (Build from source)

You need Go 1.16+ installed.

```bash
# Clone this repo
git clone https://github.com/h1code2/ssho.git
cd ssho

# Build the magic
# (We use -s -w to make the binary tiny)
go build -ldflags="-s -w" -o ssho

# Run it
./ssho
```

Then open your browser and hit `http://localhost:8080`. Welcome to the future.

## 🎮 Controls

* **Move Canvas:** Hold `Space` + `Left Click Drag` OR `Middle Mouse Button Drag`.
* **New Terminal:** Click the giant floating **`+`** button (you can't miss it).
* **Move Window:** Drag the header.
* **Ghost Mode:** Windows become translucent when dragged, because it looks cool.

## 🤖 The "AI Confession"

Full transparency: **I didn't build this alone.**

This project was built with the heavy assistance of **AI (Gemini)**. I acted as the Product Manager, Lead Architect, and "Guy who complains when the CSS is off by 1 pixel," while the AI did the heavy lifting of writing the Go and JS code.

It turns out, if you yell at an LLM enough about "Pixel Perfect Borders" and "Negative Spread Shadows," it eventually creates art.

## ❤️ Acknowledgements & Inspiration

Huge, massive shoutout to the creators of **[sshx.io](https://sshx.io)**.

When I saw `sshx`, I was blown away by its UI/UX. It's arguably the best terminal collaboration tool out there. **ssho** is my attempt to bring that specific "infinite canvas" feeling to a **local** single-player experience.

This project is a love letter to their design. If you need to collaborate with others, go use **sshx** immediately. If you want to play alone on an infinite board, use **ssho**.

## 📄 License

MIT. Do whatever you want with it.

-----

*Happy Hacking\!*