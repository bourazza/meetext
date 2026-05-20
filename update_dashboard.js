const fs = require('fs');

let code = fs.readFileSync('apps/web/app/(app)/dashboard/page.tsx', 'utf8');

// 1. Fix decisions mapping
code = code.replace(
  /{results\.decisions\.map\(\(decision, i\) => \([\s\S]*?\}\)\)}/,
  `{results.decisions.length === 0 ? (
                      <p className="text-xs text-zinc-500 italic">No decisions recorded.</p>
                    ) : (
                      results.decisions.map((decision, i) => (
                        <div key={i} className="group relative flex items-start gap-3 border border-zinc-100 rounded-lg p-3 hover:bg-zinc-50/50">
                          <div className="mt-0.5 rounded-full bg-emerald-50 p-1 text-emerald-600 border border-emerald-100">
                            <Check className="h-3.5 w-3.5" />
                          </div>
                          <div className="flex-1">
                            <p className="text-xs text-zinc-650 font-medium leading-relaxed">{decision.decision}</p>
                            <span className="inline-block mt-2 rounded bg-zinc-100 px-1.5 py-0.5 text-[9px] font-bold text-zinc-550">
                              By: {decision.made_by || "Team"}
                            </span>
                            <span className={\`inline-block mt-2 ml-2 rounded px-1.5 py-0.5 text-[9px] font-bold \${decision.confidence_score > 0.8 ? 'bg-emerald-100 text-emerald-700 border border-emerald-200' : 'bg-amber-100 text-amber-700 border border-amber-200'}\`}>
                              Confidence: {(decision.confidence_score * 100).toFixed(0)}%
                            </span>
                          </div>
                          <div className="opacity-0 group-hover:opacity-100 transition-opacity">
                            <button 
                              onClick={() => startEditing('decisions', decision.decision, i, 'decision')}
                              className="p-1 rounded text-zinc-400 hover:text-zinc-950"
                            >
                              <Edit2 className="h-3 w-3" />
                            </button>
                          </div>
                        </div>
                      ))
                    )}`
);

// 2. Fix risks mapping
code = code.replace(
  /{results\.risks\.map\(\(risk, i\) => \([\s\S]*?\}\)\)}/,
  `{results.risks.length === 0 ? (
                      <p className="text-xs text-zinc-500 italic">No risks recorded.</p>
                    ) : (
                      results.risks.map((risk, i) => (
                        <div key={i} className="group relative border border-zinc-100 rounded-lg p-3 hover:bg-zinc-50/50 transition">
                          <div className="flex items-start justify-between gap-4">
                            <div>
                              <div className="flex items-center gap-2">
                                <span className="rounded bg-rose-50 border border-rose-100 px-1.5 py-0.5 text-[9px] font-bold text-rose-600 uppercase">
                                  Severity: {risk.severity || "Unknown"}
                                </span>
                                <span className={\`rounded px-1.5 py-0.5 text-[9px] font-bold uppercase border \${risk.confidence_score > 0.8 ? 'bg-emerald-50 border-emerald-100 text-emerald-600' : 'bg-amber-50 border-amber-100 text-amber-600'}\`}>
                                  Conf: {(risk.confidence_score * 100).toFixed(0)}%
                                </span>
                              </div>
                              <p className="text-xs font-bold text-zinc-900 mt-1.5">{risk.risk}</p>
                              <div className="mt-2.5 rounded bg-zinc-100 px-3 py-2 text-[10px] text-zinc-650 border border-zinc-200">
                                <span className="font-bold text-zinc-800">Mitigation:</span> {risk.mitigation || "None provided"}
                              </div>
                            </div>
                            
                            <div className="opacity-0 group-hover:opacity-100 transition-opacity">
                              <button 
                                onClick={() => startEditing('risks', risk.risk, i, 'risk')}
                                className="p-1 rounded text-zinc-400 hover:text-zinc-950"
                              >
                                <Edit2 className="h-3 w-3" />
                              </button>
                            </div>
                          </div>
                        </div>
                      ))
                    )}`
);

// 3. Fix project documentation mapping
code = code.replace(
  /results\.project_documentation/g,
  'results.project_documentation_markdown'
);

// 4. Update collapsed cards states to include blockers and action_items
code = code.replace(
  /client_requests: false,/g,
  'blockers: false,\n    action_items: false,'
);

// Add Blockers card after Risks card
const blockersHtml = `
              {/* Card: Blockers */}
              <div className="rounded-xl border border-zinc-200/80 bg-white shadow-sm overflow-hidden">
                <div className="flex items-center justify-between border-b border-zinc-100 bg-zinc-50/50 px-5 py-4">
                  <div className="flex items-center gap-2.5">
                    <AlertCircle className="h-4 w-4 text-rose-600" />
                    <h3 className="font-bold text-zinc-900 text-sm">Active Blockers</h3>
                  </div>
                  <button 
                    onClick={() => toggleCollapse('blockers')}
                    className="p-1 rounded text-zinc-400 hover:text-zinc-950 hover:bg-zinc-100 transition"
                  >
                    {collapsedCards.blockers ? <ChevronDown className="h-4 w-4" /> : <ChevronUp className="h-4 w-4" />}
                  </button>
                </div>

                {!collapsedCards.blockers && (
                  <div className="p-5 space-y-4">
                    {results.blockers.length === 0 ? (
                      <p className="text-xs text-zinc-500 italic">No active blockers recorded.</p>
                    ) : (
                      results.blockers.map((blocker, i) => (
                        <div key={i} className="group relative border border-zinc-100 rounded-lg p-3 hover:bg-zinc-50/50 transition">
                          <div className="flex items-start justify-between gap-4">
                            <div>
                              <p className="text-xs font-bold text-zinc-900">{blocker.description}</p>
                              <span className={\`inline-block mt-2.5 rounded px-1.5 py-0.5 text-[9px] font-bold uppercase border \${blocker.confidence_score > 0.8 ? 'bg-emerald-50 border-emerald-100 text-emerald-600' : 'bg-amber-50 border-amber-100 text-amber-600'}\`}>
                                Conf: {(blocker.confidence_score * 100).toFixed(0)}%
                              </span>
                            </div>
                          </div>
                        </div>
                      ))
                    )}
                  </div>
                )}
              </div>
`;

// Add Action Items card after Blockers
const actionItemsHtml = `
              {/* Card: Action Items */}
              <div className="rounded-xl border border-zinc-200/80 bg-white shadow-sm overflow-hidden">
                <div className="flex items-center justify-between border-b border-zinc-100 bg-zinc-50/50 px-5 py-4">
                  <div className="flex items-center gap-2.5">
                    <CheckCircle2 className="h-4 w-4 text-indigo-500" />
                    <h3 className="font-bold text-zinc-900 text-sm">Action Items</h3>
                  </div>
                  <button 
                    onClick={() => toggleCollapse('action_items')}
                    className="p-1 rounded text-zinc-400 hover:text-zinc-950 hover:bg-zinc-100 transition"
                  >
                    {collapsedCards.action_items ? <ChevronDown className="h-4 w-4" /> : <ChevronUp className="h-4 w-4" />}
                  </button>
                </div>

                {!collapsedCards.action_items && (
                  <div className="p-5 space-y-4">
                    {results.action_items.length === 0 ? (
                      <p className="text-xs text-zinc-500 italic">No action items recorded.</p>
                    ) : (
                      results.action_items.map((item, i) => (
                        <div key={i} className="group relative border border-zinc-100 rounded-lg p-3 hover:bg-zinc-50/50 transition">
                          <div className="flex items-start justify-between gap-4">
                            <div>
                              <p className="text-xs font-bold text-zinc-900">{item.description}</p>
                              <div className="flex gap-2 mt-2">
                                <span className="rounded bg-zinc-100 px-1.5 py-0.5 text-[9px] font-bold text-zinc-600 uppercase border border-zinc-200">
                                  Owner: {item.owner || "Unassigned"}
                                </span>
                                <span className={\`rounded px-1.5 py-0.5 text-[9px] font-bold uppercase border \${item.confidence_score > 0.8 ? 'bg-emerald-50 border-emerald-100 text-emerald-600' : 'bg-amber-50 border-amber-100 text-amber-600'}\`}>
                                  Conf: {(item.confidence_score * 100).toFixed(0)}%
                                </span>
                              </div>
                            </div>
                          </div>
                        </div>
                      ))
                    )}
                  </div>
                )}
              </div>
`;

code = code.replace(
  /\{\/\* Card 6: Project Documentation \*\/\}/,
  blockersHtml + '\n' + actionItemsHtml + '\n              {/* Card 6: Project Documentation */}'
);

fs.writeFileSync('apps/web/app/(app)/dashboard/page.tsx', code);
console.log('Dashboard updated successfully!');
