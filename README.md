# webgo

Desktop application with Web-based UI written in Golang!


Put chromium unpacked archive to **assets** so that the chromium executable would be located in **assets/chromium/VERSION_NUM/Chrome-bin/chrome** (can be taken from: https://github.com/Hibbiki/chromium-win64/releases/download/v99.0.4844.74-r961656/chrome.sync.7z)

Run **go generate ./...**

Run **npm install** from **web-ui** folder

Run **npm run dist** from **web-ui** folder

Run **go build** if you want to build with console debug window or run **go build -ldflags -H=windowsgui** if you want to get clean GUI app without the console window

![chrome_Le523Vz57m](https://user-images.githubusercontent.com/72680690/159305077-f7887c24-7f37-485a-a115-0378437ac206.gif)
