const fs = require('fs');

let code = fs.readFileSync('apps/web/app/(app)/dashboard/page.tsx', 'utf8');

// Fix tasks
code = code.replace(
  /\{results\.tasks\.map\(\(task, i\) => \(/,
  `{(!results.tasks || results.tasks.length === 0) ? (
                      <p className="text-xs text-zinc-500 italic">No tasks recorded.</p>
                    ) : (results.tasks.map((task, i) => (`
);
code = code.replace(
  /        \}\)\)\}\n                  <\/div>/g,
  `        }))}\n                  </div>`
); // Note: We need to make sure the closing parenthesis is properly closed. 

// Actually, regex replacement for just the maps:
// The easiest way is to wrap the maps or use `(results.X || []).map` 
// But we want the empty state. Let's do string replacement for each one carefully.

// Let's rewrite the script to do safer replace.
