"use client";

import { useState } from "react";
import { motion } from "framer-motion";
import SimilarIncidentPanel from "./components/SimilarIncidentPanel";

type SimilarIncident = {
  id: number;
  title: string;
  summary: string;
  similarity: number;
};

type ExtractedIncident = {
  title: string;
  summary: string;
  symptoms: string[];
  evidence: string[];
  root_cause: string[];
  resolution: string[];
  prevention: string[];
  commands_used: string[];
  tags: string[];
  severity: string;
  environment: string;
  services_affected: string[];
  lessons_learned: string;
  raw_notes: string;
};

type Incident = ExtractedIncident & {
  id: number;
  created_at: string;
  updated_at: string;
};

export default function Home() {
  const [rawNotes, setRawNotes] = useState("");
  const [analyzing, setAnalyzing] = useState(false);
  const [extracted, setExtracted] = useState<ExtractedIncident | null>(null);
  const [saving, setSaving] = useState(false);
  const [message, setMessage] = useState("");
  const [incidentId, setIncidentId] = useState<number | null>(null);
  const [similarIncidents, setSimilarIncidents] = useState<SimilarIncident[]>([]);

  const API = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

  const handleAnalyze = async () => {
    if (!rawNotes.trim()) {
      setMessage("Please paste incident notes first");
      return;
    }
    setAnalyzing(true);
    setMessage("");
    setExtracted(null);
    setIncidentId(null);
    setSimilarIncidents([]);
    try {
      // Parallel calls: find similar incidents + extract new incident
      const [similarRes, extractRes] = await Promise.all([
        fetch(`${API}/incidents/similar`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ query: rawNotes }),
        }),
        fetch(`${API}/incident/extract`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ raw_notes: rawNotes }),
        }),
      ]);

      let simCount = 0;
      if (similarRes.ok) {
        const simData: SimilarIncident[] = await similarRes.json();
        setSimilarIncidents(simData);
        simCount = simData.length;
      }

      if (!extractRes.ok) {
        throw new Error(`Analysis failed: ${extractRes.status}`);
      }

      const data: ExtractedIncident = await extractRes.json();
      setExtracted(data);
      setMessage(
        simCount > 0
          ? `Found ${simCount} similar incident${simCount > 1 ? "s" : ""} ✅`
          : "Incident analyzed successfully ✅"
      );
    } catch (err: any) {
      console.error(err);
      setMessage(`Analysis failed: ${err.message}`);
    } finally {
      setAnalyzing(false);
    }
  };

  const handleSave = async () => {
    if (!extracted) {
      setMessage("No extracted incident to save");
      return;
    }
    setSaving(true);
    setMessage("");
    try {
      // If we already have an ID (from previous save), we update; otherwise create
      let res: Response;
      if (incidentId !== null) {
        res = await fetch(`${API}/incident/${incidentId}`, {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(extracted),
        });
      } else {
        res = await fetch(`${API}/incident`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(extracted),
        });
      }

      if (!res.ok) {
        throw new Error(`Save failed: ${res.status}`);
      }

      const data: Incident = await res.json();
      setIncidentId(data.id);
      setMessage(`Incident saved successfully ✅ (ID: ${data.id})`);
    } catch (err: any) {
      console.error(err);
      setMessage(`Save failed: ${err.message}`);
    } finally {
      setSaving(false);
    }
  };

  // Helper to convert array to newline-separated string for textarea
  const arrayToText = (arr: string[] | null | undefined): string => (arr || []).join("\n");
  // Helper to convert newline-separated string to array
  const textToArray = (text: string): string[] =>
    text
      .split("\n")
      .map((line) => line.trim())
      .filter(Boolean);

  return (
    <div className="min-h-screen bg-gradient-to-b from-[#0f172a] to-[#1e293b] text-white p-6">
      <div className="max-w-4xl mx-auto">
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          className="mb-8"
        >
          <h1 className="text-4xl font-bold text-center">
            DevOps Memory Assistant 🚀
          </h1>
          <p className="text-gray-300 mt-2 text-center">
            Transform raw troubleshooting notes into structured knowledge
          </p>
        </motion.div>

        {/* Message */}
        {message && (
          <div className={`mb-4 px-4 py-2 rounded-lg bg-black/50 text-sm text-center ${message.includes("✅") ? "text-green-400" : "text-red-400"}`}>
            {message}
          </div>
        )}

        {/* Step 1: Paste Incident Notes */}
        <motion.div
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          className="mb-6"
        >
          <div className="space-x-3 mb-2 text-sm font-medium">
            <span>Step 1:</span>
            <span>Paste Incident Notes</span>
          </div>
          <textarea
            value={rawNotes}
            onChange={(e) => setRawNotes(e.target.value)}
            placeholder="Paste your raw troubleshooting notes here...\n\nExample:\nPods entered CrashLoopBackOff after secret rotation. kubectl describe pod showed secret mount failures. Secret was accidentally deleted. Recreated secret and restarted deployment. Service recovered."
            rows={8}
            className="w-full px-4 py-3 rounded-xl bg-black/50 border border-white/10 focus:outline-none focus:ring-2 focus:ring-pink-500 text-white"
          />
        </motion.div>

        {/* Step 2: Analyze Button */}
        <motion.div
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          className="mb-6"
        >
          <div className="flex justify-center">
            <button
              onClick={handleAnalyze}
              disabled={analyzing || !rawNotes.trim()}
              className={`px-8 py-3 rounded-xl font-semibold transition-all duration-200 text-white ${
                analyzing
                  ? "bg-gray-500 cursor-not-allowed"
                  : "bg-pink-600 hover:bg-pink-700"
              }`}
            >
              {analyzing ? "Analyzing..." : "Analyze Incident"}
            </button>
          </div>
        </motion.div>

        {/* Similar Incidents Panel */}
        <SimilarIncidentPanel
          incidents={similarIncidents}
          onSelect={(inc) => {
            window.alert(`Incident #${inc.id}\n\nTitle: ${inc.title}\n\nSummary: ${inc.summary}\n\nSimilarity: ${Math.round(inc.similarity * 100)}%`);
          }}
        />

        {/* Step 3: Show and Edit Extracted Data */}
        {extracted && (
          <motion.div
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            className="space-y-6"
          >
            <div className="space-x-3 mb-2 text-sm font-medium">
              <span>Step 3:</span>
              <span>Review and Edit Extracted Incident</span>
            </div>

            {/* Form fields */}
            <div className="space-y-4">
              {/* Title */}
              <div>
                <label className="block mb-1 text-sm font-medium">Title</label>
                <input
                  value={extracted.title}
                  onChange={(e) =>
                    setExtracted((prev) => prev ? { ...prev, title: e.target.value } : null)
                  }
                  className="w-full px-4 py-2 rounded-xl bg-black/50 border border-white/10 focus:outline-none focus:ring-2 focus:ring-pink-500 text-white"
                />
              </div>

              {/* Summary */}
              <div>
                <label className="block mb-1 text-sm font-medium">Summary</label>
                <textarea
                  value={extracted.summary}
                  onChange={(e) =>
                    setExtracted((prev) => prev ? { ...prev, summary: e.target.value } : null)
                  }
                  rows={3}
                  className="w-full px-4 py-2 rounded-xl bg-black/50 border border-white/10 focus:outline-none focus:ring-2 focus:ring-pink-500 text-white"
                />
              </div>

              {/* Symptoms */}
              <div>
                <label className="block mb-1 text-sm font-medium">Symptoms (one per line)</label>
                <textarea
                  value={arrayToText(extracted.symptoms)}
                  onChange={(e) =>
                    setExtracted((prev) =>
                      prev
                        ? { ...prev, symptoms: textToArray(e.target.value) }
                        : null
                    )
                  }
                  rows={3}
                  className="w-full px-4 py-2 rounded-xl bg-black/50 border border-white/10 focus:outline-none focus:ring-2 focus:ring-pink-500 text-white"
                />
              </div>

              {/* Evidence */}
              <div>
                <label className="block mb-1 text-sm font-medium">Evidence (one per line)</label>
                <textarea
                  value={arrayToText(extracted.evidence)}
                  onChange={(e) =>
                    setExtracted((prev) =>
                      prev ? { ...prev, evidence: textToArray(e.target.value) } : null
                    )
                  }
                  rows={3}
                  className="w-full px-4 py-2 rounded-xl bg-black/50 border border-white/10 focus:outline-none focus:ring-2 focus:ring-pink-500 text-white"
                />
              </div>

              {/* Root Cause */}
              <div>
                <label className="block mb-1 text-sm font-medium">Root Cause (one per line)</label>
                <textarea
                  value={arrayToText(extracted.root_cause)}
                  onChange={(e) =>
                    setExtracted((prev) =>
                      prev
                        ? { ...prev, root_cause: textToArray(e.target.value) }
                        : null
                    )
                  }
                  rows={3}
                  className="w-full px-4 py-2 rounded-xl bg-black/50 border border-white/10 focus:outline-none focus:ring-2 focus:ring-pink-500 text-white"
                />
              </div>

              {/* Resolution */}
              <div>
                <label className="block mb-1 text-sm font-medium">Resolution (one per line)</label>
                <textarea
                  value={arrayToText(extracted.resolution)}
                  onChange={(e) =>
                    setExtracted((prev) =>
                      prev
                        ? { ...prev, resolution: textToArray(e.target.value) }
                        : null
                    )
                  }
                  rows={3}
                  className="w-full px-4 py-2 rounded-xl bg-black/50 border border-white/10 focus:outline-none focus:ring-2 focus:ring-pink-500 text-white"
                />
              </div>

              {/* Commands Used */}
              <div>
                <label className="block mb-1 text-sm font-medium">Commands Used (one per line)</label>
                <textarea
                  value={arrayToText(extracted.commands_used)}
                  onChange={(e) =>
                    setExtracted((prev) =>
                      prev
                        ? { ...prev, commands_used: textToArray(e.target.value) }
                        : null
                    )
                  }
                  rows={3}
                  className="w-full px-4 py-2 rounded-xl bg-black/50 border border-white/10 focus:outline-none focus:ring-2 focus:ring-pink-500 text-white"
                />
              </div>

              {/* Prevention */}
              <div>
                <label className="block mb-1 text-sm font-medium">Prevention (one per line)</label>
                <textarea
                  value={arrayToText(extracted.prevention)}
                  onChange={(e) =>
                    setExtracted((prev) =>
                      prev
                        ? { ...prev, prevention: textToArray(e.target.value) }
                        : null
                    )
                  }
                  rows={3}
                  className="w-full px-4 py-2 rounded-xl bg-black/50 border border-white/10 focus:outline-none focus:ring-2 focus:ring-pink-500 text-white"
                />
              </div>

              {/* Tags */}
              <div>
                <label className="block mb-1 text-sm font-medium">Tags (comma-separated)</label>
                <input
                  value={extracted.tags.join(", ")}
                  onChange={(e) =>
                    setExtracted((prev) =>
                      prev
                        ? {
                            ...prev,
                            tags: e.target.value
                              .split(",")
                              .map((tag) => tag.trim())
                              .filter(Boolean),
                          }
                        : null
                    )
                  }
                  className="w-full px-4 py-2 rounded-xl bg-black/50 border border-white/10 focus:outline-none focus:ring-2 focus:ring-pink-500 text-white"
                />
              </div>

              {/* Severity */}
              <div>
                <label className="block mb-1 text-sm font-medium">Severity</label>
                <select
                  value={extracted.severity}
                  onChange={(e) =>
                    setExtracted((prev) => prev ? { ...prev, severity: e.target.value } : null)
                  }
                  className="w-full px-4 py-2 rounded-xl bg-black/50 border border-white/10 focus:outline-none focus:ring-2 focus:ring-pink-500 text-white"
                >
                  <option value="low">Low</option>
                  <option value="medium">Medium</option>
                  <option value="high">High</option>
                  <option value="critical">Critical</option>
                </select>
              </div>

              {/* Environment */}
              <div>
                <label className="block mb-1 text-sm font-medium">Environment</label>
                <input
                  value={extracted.environment}
                  onChange={(e) =>
                    setExtracted((prev) =>
                      prev ? { ...prev, environment: e.target.value } : null
                    )
                  }
                  className="w-full px-4 py-2 rounded-xl bg-black/50 border border-white/10 focus:outline-none focus:ring-2 focus:ring-pink-500 text-white"
                />
              </div>

              {/* Services Affected */}
              <div>
                <label className="block mb-1 text-sm font-medium">Services Affected (one per line)</label>
                <textarea
                  value={arrayToText(extracted.services_affected)}
                  onChange={(e) =>
                    setExtracted((prev) =>
                      prev
                        ? {
                            ...prev,
                            services_affected: textToArray(e.target.value),
                          }
                        : null
                    )
                  }
                  rows={2}
                  className="w-full px-4 py-2 rounded-xl bg-black/50 border border-white/10 focus:outline-none focus:ring-2 focus:ring-pink-500 text-white"
                />
              </div>

              {/* Lessons Learned */}
              <div>
                <label className="block mb-1 text-sm font-medium">Lessons Learned</label>
                <textarea
                  value={extracted.lessons_learned}
                  onChange={(e) =>
                    setExtracted((prev) =>
                      prev ? { ...prev, lessons_learned: e.target.value } : null
                    )
                  }
                  rows={3}
                  className="w-full px-4 py-2 rounded-xl bg-black/50 border border-white/10 focus:outline-none focus:ring-2 focus:ring-pink-500 text-white"
                />
              </div>
            </div>

            {/* Step 4: Save Button */}
            <motion.div
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              className="mt-6"
            >
              <div className="flex justify-center">
                <button
                  onClick={handleSave}
                  disabled={saving}
                  className={`px-8 py-3 rounded-xl font-semibold transition-all duration-200 text-white ${
                    saving
                      ? "bg-gray-500 cursor-not-allowed"
                      : "bg-purple-600 hover:bg-purple-700"
                  }`}
                >
                  {saving ? "Saving..." : "Save Incident"}
                </button>
              </div>
            </motion.div>
          </motion.div>
        )}

        {/* Optional: Show saved incident ID */}
        {incidentId && (
          <motion.div
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            className="mt-6 px-4 py-2 rounded-lg bg-black/50 text-center text-sm"
          >
            Incident saved with ID: <span className="font-semibold text-pink-300">{incidentId}</span>
          </motion.div>
        )}
      </div>
    </div>
  );
}