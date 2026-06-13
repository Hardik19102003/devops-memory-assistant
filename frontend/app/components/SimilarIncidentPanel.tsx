"use client";

import { motion } from "framer-motion";

export type SimilarIncident = {
  id: number;
  title: string;
  summary: string;
  similarity: number;
};

type Props = {
  incidents: SimilarIncident[];
  onSelect?: (incident: SimilarIncident) => void;
};

export default function SimilarIncidentPanel({ incidents, onSelect }: Props) {
  if (!incidents || incidents.length === 0) return null;

  return (
    <motion.div
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      className="mb-6"
    >
      <div className="flex items-center gap-2 mb-3">
        <span className="text-sm font-semibold text-pink-400">
          🔍 Similar Incidents Found
        </span>
        <span className="text-xs text-gray-400">
          ({incidents.length} match{incidents.length > 1 ? "es" : ""})
        </span>
      </div>

      <div className="grid gap-3">
        {incidents.map((inc, idx) => {
          const pct = Math.round(inc.similarity * 100);
          const badgeColor =
            pct >= 90
              ? "bg-green-500/20 text-green-400 border-green-500/30"
              : pct >= 75
              ? "bg-yellow-500/20 text-yellow-400 border-yellow-500/30"
              : "bg-orange-500/20 text-orange-400 border-orange-500/30";

          return (
            <motion.div
              key={inc.id}
              initial={{ opacity: 0, x: -10 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: idx * 0.08 }}
              onClick={() => onSelect?.(inc)}
              className="cursor-pointer group relative rounded-xl border border-white/10 bg-black/40 p-4 hover:bg-black/60 hover:border-pink-500/40 transition-all duration-200"
            >
              <div className="flex items-start justify-between gap-3">
                <div className="flex-1 min-w-0">
                  <h3 className="text-sm font-semibold text-white truncate group-hover:text-pink-300 transition-colors">
                    {inc.title}
                  </h3>
                  <p className="text-xs text-gray-400 mt-1 line-clamp-2">
                    {inc.summary || "No summary available"}
                  </p>
                </div>
                <div
                  className={`shrink-0 inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-bold border ${badgeColor}`}
                >
                  <span>{pct}%</span>
                  <span className="opacity-70">Match</span>
                </div>
              </div>
            </motion.div>
          );
        })}
      </div>
    </motion.div>
  );
}
