"use client";

import { useState } from "react";
import { motion } from "framer-motion";

type Issue = {
  id: number;
  error: string;
  causes: string[];
  fixes: string[];
  debug_steps: string[];
  tags: string[];
  document: string;
  created_at: string;
};

export default function Home() {
  const [error, setError] = useState("");
  const [causes, setCauses] = useState("");
  const [fixes, setFixes] = useState("");
  const [debugSteps, setDebugSteps] = useState("");
  const [tags, setTags] = useState("");
  const [document, setDocument] = useState("");

  const [results, setResults] = useState<Issue[]>([]);
  const [suggestions, setSuggestions] = useState<Issue[]>([]);

  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState("");
  const [searched, setSearched] = useState(false);

  const API = process.env.NEXT_PUBLIC_API_URL;

  // SEARCH
  const handleSearch = async () => {
    setLoading(true);
    setMessage("");
    setSearched(true);

    try {
      const res = await fetch(
        `${API}/search?error=${encodeURIComponent(error)}`
      );

      const data = await res.json();

      setResults(data.results || []);
    } catch {
      setMessage("Search failed ❌");
    }

    setLoading(false);
  };

  // SAVE
  const handleSave = async () => {
    setLoading(true);
    setMessage("");

    try {
      const res = await fetch(`${API}/issue`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: "Bearer devops-secret-key",
        },
        body: JSON.stringify({
          error,

          causes: causes
            .split("\n")
            .map((x) => x.trim())
            .filter(Boolean),

          fixes: fixes
            .split("\n")
            .map((x) => x.trim())
            .filter(Boolean),

          debug_steps: debugSteps
            .split("\n")
            .map((x) => x.trim())
            .filter(Boolean),

          tags: tags
            .split(",")
            .map((x) => x.trim())
            .filter(Boolean),

          document,
        }),
      });

      const data = await res.json();

      setMessage(data.message || "Saved ✅");

      setError("");
      setCauses("");
      setFixes("");
      setDebugSteps("");
      setTags("");
      setDocument("");
    } catch {
      setMessage("Save failed ❌");
    }

    setLoading(false);
  };

  // LIVE SEARCH SUGGESTIONS
  const handleInputChange = async (value: string) => {
    setError(value);

    if (value.length > 2) {
      try {
        const res = await fetch(
          `${API}/search?error=${encodeURIComponent(value)}`
        );

        const data = await res.json();

        setSuggestions(data.results || []);
      } catch {
        setSuggestions([]);
      }
    } else {
      setSuggestions([]);
    }
  };

  // DELETE
  const handleDelete = async (id: number) => {
    try {
      await fetch(`${API}/delete?id=${id}`, {
        method: "DELETE",
        headers: {
          Authorization: "Bearer devops-secret-key",
        },
      });

      setResults((prev) =>
        prev.filter((item) => item.id !== id)
      );

      setMessage("Deleted successfully ✅");
    } catch {
      setMessage("Delete failed ❌");
    }
  };

  return (
    <div className="min-h-screen bg-[#0f172a] text-white p-6">
      <div className="max-w-7xl mx-auto">

        {/* TITLE */}
        <motion.div
          initial={{ opacity: 0, y: -15 }}
          animate={{ opacity: 1, y: 0 }}
          className="mb-8"
        >
          <h1 className="text-5xl font-bold">
            DevOps Memory Assistant 🚀
          </h1>

          <p className="text-gray-400 mt-2">
            Store production incidents, debugging workflows,
            operational notes, and resolutions.
          </p>
        </motion.div>

        {/* MAIN GRID */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">

          {/* LEFT PANEL */}
          <div className="bg-white/5 border border-white/10 rounded-2xl p-6 space-y-4">

            {/* INCIDENT TITLE */}
            <div>
              <label className="text-sm text-gray-300">
                Incident / Error
              </label>

              <input
                value={error}
                onChange={(e) =>
                  handleInputChange(e.target.value)
                }
                placeholder="CrashLoopBackOff"
                className="w-full mt-2 p-3 rounded-xl bg-black/30 border border-white/10"
              />

              {/* SUGGESTIONS */}
              {suggestions.length > 0 && (
                <div className="mt-2 bg-black/80 rounded-xl border border-white/10 overflow-hidden">
                  {suggestions.map((item) => (
                    <div
                      key={item.id}
                      onClick={() => {
                        setError(item.error);
                        setSuggestions([]);
                      }}
                      className="p-3 hover:bg-white/10 cursor-pointer border-b border-white/5"
                    >
                      {item.error}
                    </div>
                  ))}
                </div>
              )}
            </div>

            {/* CAUSES */}
            <textarea
              value={causes}
              onChange={(e) => setCauses(e.target.value)}
              placeholder="Causes (one per line)"
              rows={4}
              className="w-full p-3 rounded-xl bg-black/30 border border-white/10"
            />

            {/* FIXES */}
            <textarea
              value={fixes}
              onChange={(e) => setFixes(e.target.value)}
              placeholder="Fixes (one per line)"
              rows={4}
              className="w-full p-3 rounded-xl bg-black/30 border border-white/10"
            />

            {/* DEBUG STEPS */}
            <textarea
              value={debugSteps}
              onChange={(e) => setDebugSteps(e.target.value)}
              placeholder="Debug steps (one per line)"
              rows={5}
              className="w-full p-3 rounded-xl bg-black/30 border border-white/10"
            />

            {/* TAGS */}
            <input
              value={tags}
              onChange={(e) => setTags(e.target.value)}
              placeholder="kubernetes,docker,network"
              className="w-full p-3 rounded-xl bg-black/30 border border-white/10"
            />

            {/* DOCUMENT */}
            <div>
              <label className="text-sm text-gray-300">
                Incident Documentation / Runbook
              </label>

              <textarea
                value={document}
                onChange={(e) => setDocument(e.target.value)}
                placeholder={`Paste full incident notes here...

Example:
- What happened
- Root cause
- Investigation timeline
- Commands used
- Recovery steps
- Permanent fix
- Lessons learned`}
                rows={18}
                className="w-full mt-2 p-4 rounded-xl bg-black/30 border border-white/10 font-mono text-sm"
              />
            </div>

            {/* BUTTONS */}
            <div className="flex gap-4">

              <button
                onClick={handleSearch}
                disabled={loading}
                className="flex-1 bg-pink-600 hover:bg-pink-700 py-3 rounded-xl font-semibold"
              >
                {loading ? "Searching..." : "Search"}
              </button>

              <button
                onClick={handleSave}
                disabled={loading}
                className="flex-1 bg-purple-600 hover:bg-purple-700 py-3 rounded-xl font-semibold"
              >
                {loading ? "Saving..." : "Save"}
              </button>
            </div>

            {/* MESSAGE */}
            {message && (
              <div className="bg-black/30 border border-white/10 p-3 rounded-xl text-sm">
                {message}
              </div>
            )}
          </div>

          {/* RIGHT PANEL */}
          <div className="space-y-5">

            {searched && results.length === 0 && (
              <div className="bg-white/5 border border-white/10 rounded-2xl p-6 text-gray-400">
                No incidents found 👀
              </div>
            )}

            {results.map((item) => (
              <motion.div
                key={item.id}
                initial={{ opacity: 0, y: 15 }}
                animate={{ opacity: 1, y: 0 }}
                className="bg-white/5 border border-white/10 rounded-2xl p-6"
              >
                {/* TITLE */}
                <h2 className="text-2xl font-bold text-pink-300">
                  🚨 {item.error}
                </h2>

                {/* TAGS */}
                <div className="flex flex-wrap gap-2 mt-3">
                  {item.tags?.map((tag, idx) => (
                    <span
                      key={idx}
                      className="bg-pink-500/20 text-pink-200 px-3 py-1 rounded-full text-xs"
                    >
                      #{tag}
                    </span>
                  ))}
                </div>

                {/* CAUSES */}
                <div className="mt-5">
                  <p className="font-semibold text-yellow-300">
                    📌 Causes
                  </p>

                  <ul className="list-disc ml-5 mt-2 text-gray-300 space-y-1">
                    {item.causes?.map((cause, idx) => (
                      <li key={idx}>{cause}</li>
                    ))}
                  </ul>
                </div>

                {/* FIXES */}
                <div className="mt-5">
                  <p className="font-semibold text-green-300">
                    ✅ Fixes
                  </p>

                  <ul className="list-disc ml-5 mt-2 text-gray-300 space-y-1">
                    {item.fixes?.map((fix, idx) => (
                      <li key={idx}>{fix}</li>
                    ))}
                  </ul>
                </div>

                {/* DEBUG */}
                <div className="mt-5">
                  <p className="font-semibold text-blue-300">
                    🛠 Debug Steps
                  </p>

                  <ul className="list-disc ml-5 mt-2 text-gray-300 space-y-1">
                    {item.debug_steps?.map((step, idx) => (
                      <li key={idx}>{step}</li>
                    ))}
                  </ul>
                </div>

                {/* DOCUMENT */}
                {item.document && (
                  <div className="mt-6">
                    <p className="font-semibold text-purple-300 mb-3">
                      📚 Incident Notes
                    </p>

                    <div className="bg-black/40 border border-white/10 rounded-xl p-4 whitespace-pre-wrap text-sm text-gray-300 font-mono max-h-[400px] overflow-y-auto">
                      {item.document}
                    </div>
                  </div>
                )}

                {/* FOOTER */}
                <div className="flex items-center justify-between mt-6">

                  <p className="text-xs text-gray-500">
                    🕒{" "}
                    {new Date(
                      item.created_at
                    ).toLocaleString()}
                  </p>

                  <button
                    onClick={() =>
                      handleDelete(item.id)
                    }
                    className="bg-red-500 hover:bg-red-600 px-3 py-1 rounded-lg text-sm"
                  >
                    🗑 Delete
                  </button>
                </div>
              </motion.div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}