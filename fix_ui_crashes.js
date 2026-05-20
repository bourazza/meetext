const fs = require('fs');

let code = fs.readFileSync('apps/web/app/(app)/dashboard/page.tsx', 'utf8');

// 1. Tasks
code = code.replace(
  /\{results\.tasks\.map\(\(task, i\) => \(/,
  `{(!results.tasks || results.tasks.length === 0) ? (
                      <p className="text-xs text-zinc-500 italic">No tasks recorded.</p>
                    ) : (
                      results.tasks.map((task, i) => (`
);
// Find the end of tasks map (lines 880-885)
code = code.replace(
  /                            <\/button>\n                          <\/div>\n                        <\/div>\n                      <\/div>\n                    \)\)\}/,
  `                            </button>\n                          </div>\n                        </div>\n                      </div>\n                    )))}`
);

// 2. Tickets
code = code.replace(
  /\{results\.tickets\.map\(\(ticket, i\) => \(/,
  `{(!results.tickets || results.tickets.length === 0) ? (
                      <p className="text-xs text-zinc-500 italic">No tickets recorded.</p>
                    ) : (
                      results.tickets.map((ticket, i) => (`
);
code = code.replace(
  /                            <\/button>\n                          <\/div>\n                        <\/div>\n                      <\/div>\n                    \)\)\}/,
  `                            </button>\n                          </div>\n                        </div>\n                      </div>\n                    )))}`
);

// 3. Decisions
code = code.replace(
  /\{results\.decisions\.map\(\(decision, i\) => \(/,
  `{(!results.decisions || results.decisions.length === 0) ? (
                      <p className="text-xs text-zinc-500 italic">No decisions recorded.</p>
                    ) : (
                      results.decisions.map((decision, i) => (`
);
code = code.replace(
  /                            <\/button>\n                        <\/div>\n                      <\/div>\n                    \)\)\}/,
  `                            </button>\n                        </div>\n                      </div>\n                    )))}`
);

// 4. Risks
code = code.replace(
  /\{results\.risks\.map\(\(risk, i\) => \(/,
  `{(!results.risks || results.risks.length === 0) ? (
                      <p className="text-xs text-zinc-500 italic">No risks recorded.</p>
                    ) : (
                      results.risks.map((risk, i) => (`
);
code = code.replace(
  /                            <\/button>\n                          <\/div>\n                        <\/div>\n                      <\/div>\n                    \)\)\}/,
  `                            </button>\n                          </div>\n                        </div>\n                      </div>\n                    )))}`
);

// 5. Blockers (Already has length check but let's add null check)
code = code.replace(
  /\{results\.blockers\.length === 0 \? \(/,
  `{(!results.blockers || results.blockers.length === 0) ? (`
);

// 6. Action Items
code = code.replace(
  /\{results\.action_items\.length === 0 \? \(/,
  `{(!results.action_items || results.action_items.length === 0) ? (`
);

// 7. Tech Notes
code = code.replace(
  /\{results\.technical_notes\.map\(\(note, i\) => \(/,
  `{(!results.technical_notes || results.technical_notes.length === 0) ? (
                      <p className="text-xs text-zinc-500 italic">No technical notes recorded.</p>
                    ) : (
                      results.technical_notes.map((note, i) => (`
);
code = code.replace(
  /                        <\/div>\n                      <\/div>\n                    \)\)\}/,
  `                        </div>\n                      </div>\n                    )))}`
);


fs.writeFileSync('apps/web/app/(app)/dashboard/page.tsx', code);
console.log("Success! Fallbacks applied.");
