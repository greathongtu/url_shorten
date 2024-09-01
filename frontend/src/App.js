import React, { useState } from "react";
import axios from "axios";

function App() {
  const [longUrl, setLongUrl] = useState("");
  const [shortCode, setShortCode] = useState("");
  const [error, setError] = useState("");

  const formatUrl = (url) => {
    if (!/^https?:\/\//i.test(url)) {
      return `http://${url}`;
    }
    return url;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    const formattedUrl = formatUrl(longUrl);
    try {
      const response = await axios.post("http://localhost:8080/shorten", {
        long_url: formattedUrl,
      });
      setShortCode(response.data.short_code);
    } catch (error) {
      console.error("Error shortening URL:", error);
      setError("Failed to shorten URL. Please try again.");
    }
  };

  const getShortUrl = () => `http://localhost:8080/${shortCode}`;

  return (
    <div className="App">
      <h1>URL Shortener</h1>
      <form onSubmit={handleSubmit}>
        <input
          type="text"
          value={longUrl}
          onChange={(e) => setLongUrl(e.target.value)}
          placeholder="Enter long URL"
        />
        <button type="submit">Shorten</button>
      </form>
      {error && <p style={{ color: "red" }}>{error}</p>}
      {shortCode && (
        <div>
          <h2>Shortened URL:</h2>
          <p>
            <a href={getShortUrl()} target="_blank" rel="noopener noreferrer">
              {getShortUrl()}
            </a>
          </p>
          <p>This short URL will redirect to: {formatUrl(longUrl)}</p>
        </div>
      )}
    </div>
  );
}

export default App;
