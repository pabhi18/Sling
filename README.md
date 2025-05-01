```
░██████╗██╗░░░░░██╗███╗░░██╗░██████╗░
██╔════╝██║░░░░░██║████╗░██║██╔════╝░
╚█████╗░██║░░░░░██║██╔██╗██║██║░░██╗░
░╚═══██╗██║░░░░░██║██║╚████║██║░░╚██╗
██████╔╝███████╗██║██║░╚███║╚██████╔╝
╚═════╝░╚══════╝╚═╝╚═╝░░╚══╝░╚═════╝░
```

# Sling 🔗📂

**Sling** is a simple, blazing-fast CLI tool written in Go that lets you share any folder over HTTP — right from your terminal. No external servers, no configuration — just run and sling!

## 🚀 Features

- Serve any local folder over HTTP instantly  
- Auto-detects internal IP so devices on the same network can access it  
- Simple, single-command interface  

## 🧪 Installation

Use the one-liner script below to install Sling:

```bash
curl -L https://raw.githubusercontent.com/pabhi18/sling/main/install.sh | sh
```

This will download the latest version and place it in `/usr/local/bin/sling`.

## 📂 Usage

### Serve a folder on the current network

```bash
sling serve
```

- Serves the **current directory** on port `8080`
- Accessible from other devices on the same network using your **internal IP**

### Serve a specific folder on a custom port

```bash
sling serve -path ./your-folder -p 9090
```

- Replace `./your-folder` with the path you want to serve
- Change `9090` to any port you prefer

## 💡 Example Use Cases

- Instantly share files across your devices
- Let teammates quickly download a build or logs
- Temporary local file hosting during testing or development


## 🛠️ Build from Source

```bash
git clone https://github.com/pabhi18/Sling.git
cd Sling
go build -o sling
```


## 🧾 License

MIT License


## 🤝 Contributing

Pull requests, feature ideas, and improvements are welcome!  

Made with ❤️ by [Abhinav Pratap](https://github.com/pabhi18)
