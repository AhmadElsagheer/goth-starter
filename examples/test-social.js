/**
 * BetterAuth Google Social Login Test Server
 * * Usage: 
 * 1. Ensure your BetterAuth backend is running
 * 2. Run this script: bun test-social.js
 * 3. Open http://localhost:3000 in your browser
 */

const http = require('http');
const url = require('url');

// CONFIGURATION
const PORT = 3000; // The port this test tool will run on

const htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>BetterAuth Test Client</title>
    <style>
        body { font-family: -apple-system, sans-serif; max-width: 600px; margin: 40px auto; padding: 20px; line-height: 1.6; background: #f4f4f5; color: #333; }
        .card { background: white; padding: 30px; border-radius: 12px; box-shadow: 0 4px 6px -1px rgba(0,0,0,0.1); }
        h1 { margin-top: 0; color: #111; }
        label { display: block; margin-bottom: 8px; font-weight: 600; font-size: 0.9em; }
        input { width: 100%; padding: 10px; margin-bottom: 20px; border: 1px solid #ddd; border-radius: 6px; box-sizing: border-box; }
        button { width: 100%; padding: 12px; background: #2563eb; color: white; border: none; border-radius: 6px; font-weight: bold; cursor: pointer; transition: background 0.2s; }
        button:hover { background: #1d4ed8; }
        .code-block { background: #1e293b; color: #e2e8f0; padding: 15px; border-radius: 6px; overflow-x: auto; font-family: monospace; font-size: 0.85em; margin-top: 20px; }
        .badge { display: inline-block; padding: 4px 8px; border-radius: 4px; font-size: 0.8em; font-weight: bold; }
        .badge-success { background: #dcfce7; color: #166534; }
        .badge-fail { background: #fee2e2; color: #991b1b; }
    </style>
</head>
<body>

    <div id="login-view" class="card">
        <h1>🧪 BetterAuth Tester</h1>
        <p>This tool simulates a frontend client initiating a Google Social Login.</p>
        
        <form id="loginForm">
            <label for="authUrl">Your BetterAuth Backend URL</label>
            <input type="url" id="authUrl" value="http://localhost:8787" placeholder="e.g. http://localhost:8787" required>

            <button type="submit">Login with Google</button>
        </form>
        <div id="logs" style="margin-top:20px;"></div>
    </div>

    <script>
        const form = document.getElementById('loginForm');
        const logDiv = document.getElementById('logs');

        form.addEventListener('submit', async (e) => {
            e.preventDefault();
            const baseUrl = document.getElementById('authUrl').value.replace(/\\/$/, '');
            const callbackUrl = "http://localhost:${PORT}/verified"; // This server

            logDiv.innerHTML = '<div class="code-block">Requesting auth URL...</div>';

            try {
                // 1. Call BetterAuth to get the Google Redirect URL
                const response = await fetch(baseUrl + "/api/auth/sign-in/social", {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify({
                        provider: "google",
                        callbackURL: callbackUrl 
                    }),
                    credentials: "include"
                });

                const data = await response.json();

                if (!response.ok) throw new Error(data.message || "Failed to initiate");

                logDiv.innerHTML += '<div class="code-block">Received Redirect URL:\\n' + data.url + '</div>';
                logDiv.innerHTML += '<p>Redirecting browser...</p>';

                // 2. Redirect the browser to Google
                setTimeout(() => {
                    window.location.href = data.url;
                }, 1000);

            } catch (err) {
                logDiv.innerHTML += '<div class="code-block" style="color:#f87171">Error: ' + err.message + '</div>';
            }
        });
    </script>
</body>
</html>
`;

const successTemplate = (cookies) => `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Login Verified</title>
    <style>
        body { font-family: -apple-system, sans-serif; max-width: 600px; margin: 40px auto; padding: 20px; line-height: 1.6; background: #f0fdf4; color: #333; }
        .card { background: white; padding: 30px; border-radius: 12px; box-shadow: 0 4px 6px -1px rgba(0,0,0,0.1); border: 1px solid #bbf7d0; }
        h1 { color: #166534; margin-top:0; }
        .token-box { word-break: break-all; background: #f8fafc; padding: 15px; border: 1px solid #e2e8f0; border-radius: 6px; font-family: monospace; font-size: 0.9em; }
        .btn { display: inline-block; margin-top: 20px; text-decoration: none; background: #166534; color: white; padding: 10px 20px; border-radius: 6px; font-weight: bold; }
    </style>
</head>
<body>
    <div class="card">
        <h1>✅ Login Successful!</h1>
        <p>The OAuth flow completed and redirected back to this test server.</p>
        
        <h3>Session Evidence</h3>
        <p>Because this test server is running on <code>localhost</code> (port ${PORT}) and your API is likely on <code>localhost</code> (port 8787), browsers share cookies across ports.</p>
        
        <p><strong>Cookies received by Test Server:</strong></p>
        <div class="token-box">
            ${cookies || "No cookies detected. (Check if your browser blocked them or if using different domains)"}
        </div>

        <a href="/" class="btn">Try Again</a>
    </div>
</body>
</html>
`;

const server = http.createServer((req, res) => {
    const reqUrl = url.parse(req.url, true);

    // 1. Serve the Login Page
    if (reqUrl.pathname === '/') {
        res.writeHead(200, { 'Content-Type': 'text/html' });
        res.end(htmlTemplate);
    }
    // 2. Serve the Callback Success Page
    else if (reqUrl.pathname === '/verified') {
        const cookies = req.headers.cookie;
        res.writeHead(200, { 'Content-Type': 'text/html' });
        res.end(successTemplate(cookies));
    }
    // 404
    else {
        res.writeHead(404);
        res.end('Not found');
    }
});

server.listen(PORT, () => {
    console.log("----------------------------------------------------------");
    console.log(`🚀 Test Server Running at: http://localhost:${PORT}`);
    console.log("----------------------------------------------------------");
});
