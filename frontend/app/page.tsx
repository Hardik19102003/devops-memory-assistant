"use client";

import { useState } from "react";
import { motion } from "framer-motion";

export default function Home() {
  const [error, setError] = useState("");
  const [cause, setCause] = useState("");
  const [fix, setFix] = useState("");
  const [results, setResults] = useState<any[]>([]);
  const [suggestions, setSuggestions] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState("");
  const [searched, setSearched] = useState(false);
  const [similarIssue, setSimilarIssue] = useState<any | null>(null);

  const API = process.env.NEXT_PUBLIC_API_URL;

  // 🔍 SEARCH
  const handleSearch = async () => {
    setLoading(true);
    setMessage("");
    setSearched(true);
    setResults([]);

    try {
      const res = await fetch(`${API}/search?error=${error}`);
      const data = await res.json();
      setResults(data);
    } catch {
      setMessage("Something went wrong ❌");
    }

    setLoading(false);
  };

  // 💾 SAVE
  const handleSave = async () => {
    setLoading(true);
    setMessage("");
    setSimilarIssue(null);

    try {
      const res = await fetch(`${API}/issue`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ error, cause, fix }),
      });

      const data = await res.json();

      if (data.existing) {
        setSimilarIssue(data.existing);
        setMessage("Similar issue found ⚠️");
      } else {
        setMessage("Saved successfully ✅");
        setError("");
        setCause("");
        setFix("");
      }
    } catch {
      setMessage("Failed to save ❌");
    }

    setLoading(false);
  };

  // ⚡ LIVE SUGGESTIONS
  const handleInputChange = async (value: string) => {
    setError(value);

    if (value.length > 2) {
      try {
        const res = await fetch(`${API}/search?error=${value}`);
        const data = await res.json();
        setSuggestions(data);
      } catch {
        setSuggestions([]);
      }
    } else {
      setSuggestions([]);
    }
  };

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      className="min-h-screen bg-gradient-to-br from-purple-900 via-purple-800 to-pink-700 text-white flex flex-col items-center p-6"
    >
      {/* TITLE */}
      <motion.h1
        initial={{ y: -20, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        className="text-4xl font-bold mb-6 text-center"
      >
        DevOps Memory Assistant 🚀
      </motion.h1>

      {/* FORM */}
      <motion.div
        initial={{ scale: 0.95, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        className="bg-white/10 backdrop-blur-lg p-6 rounded-2xl shadow-lg w-full max-w-md space-y-4"
      >
        {/* ERROR INPUT */}
        <div>
          <input
            value={error}
            onChange={(e) => handleInputChange(e.target.value)}
            placeholder="Error (e.g. CrashLoopBackOff)"
            className="w-full p-3 rounded-lg bg-white/90 text-black focus:ring-2 focus:ring-pink-400"
          />

          {/* 🔥 Suggestions */}
          {suggestions && suggestions.length > 0 && (
            <div className="bg-white text-black rounded-lg mt-2 shadow-lg max-h-40 overflow-y-auto">
              {suggestions.map((item, index) => (
                <div
                  key={index}
                  className="p-2 hover:bg-gray-200 cursor-pointer"
                  onClick={() => {
                    setError(item.error);
                    setSuggestions([]);
                  }}
                >
                  {item.error}
                </div>
              ))}
            </div>
          )}
        </div>

        {/* CAUSE */}
        <input
          value={cause}
          onChange={(e) => setCause(e.target.value)}
          placeholder="Cause"
          className="w-full p-3 rounded-lg bg-white/90 text-black"
        />

        {/* FIX */}
        <input
          value={fix}
          onChange={(e) => setFix(e.target.value)}
          placeholder="Fix"
          className="w-full p-3 rounded-lg bg-white/90 text-black"
        />

        {/* BUTTONS */}
        <div className="flex gap-4">
          <motion.button
            whileTap={{ scale: 0.95 }}
            onClick={handleSearch}
            disabled={loading}
            className="flex-1 bg-pink-500 py-2 rounded-lg hover:bg-pink-600 transition"
          >
            {loading ? "Searching..." : "Search"}
          </motion.button>

          <motion.button
            whileTap={{ scale: 0.95 }}
            onClick={handleSave}
            disabled={loading}
            className="flex-1 bg-purple-500 py-2 rounded-lg hover:bg-purple-600 transition"
          >
            {loading ? "Saving..." : "Save"}
          </motion.button>
        </div>

        {/* MESSAGE */}
        {message && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            className="text-center text-sm bg-black/30 p-2 rounded"
          >
            {message}
          </motion.div>
        )}

        {/* SIMILAR ISSUE WARNING */}
        {similarIssue && (
          <div className="mt-4 bg-yellow-200 text-black p-4 rounded-lg">
            <p className="font-bold">⚠️ Similar issue already exists</p>
            <p><strong>Error:</strong> {similarIssue.error}</p>
            <p><strong>Cause:</strong> {similarIssue.cause}</p>
            <p><strong>Fix:</strong> {similarIssue.fix}</p>
          </div>
        )}
      </motion.div>

      {/* RESULTS */}
      <div className="mt-8 w-full max-w-md space-y-4">
        {searched && results && results.length === 0 && !loading && (
          <p className="text-center text-gray-300">
            No results found 👀
          </p>
        )}

        {(results || []).map((item, index) => (
          <motion.div
            key={index}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: index * 0.1 }}
            className="bg-white/10 backdrop-blur-lg p-4 rounded-xl border border-white/20 hover:scale-[1.03] transition"
          >
            <p className="text-pink-300 font-semibold">{item.error}</p>
            <p className="text-sm text-gray-200">
              Cause: {item.cause}
            </p>
            <p className="text-sm text-green-300">
              Fix: {item.fix}
            </p>
          </motion.div>
        ))}
      </div>
    </motion.div>
  );
}