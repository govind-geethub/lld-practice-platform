// LLD Practice Platform - Client Application
// Supports seamless execution whether served directly from Go (port 8080),
// through external dev servers (e.g. Live Server on port 5500), or loaded directly via file://.

// Detect backend endpoint dynamically
const isLocalFile = window.location.protocol === 'file:';
const isDevServer = window.location.port && window.location.port !== '8080';
const API_BASE = (isLocalFile || isDevServer) ? 'http://localhost:8080' : '';

let selectedProblemId = null;
let problemsData = [];
let lastResult = null;

// Built-in fallback catalog to ensure problems are always visible
const FALLBACK_PROBLEMS = [
    {
        id: 1,
        title: "Parking Lot",
        difficulty: "Medium",
        description: "Design a parking lot supporting multiple floors, vehicle types (Car, Bike, Truck), hourly pricing strategies, and parking spot allocation.",
        skeleton: "class ParkingLot {\n  // TODO: Add floors, capacity, entry/exit gates\n}\n\ninterface PricingStrategy {\n  calculateCost(Ticket t): float\n}"
    },
    {
        id: 2,
        title: "Elevator System",
        difficulty: "Hard",
        description: "Design an intelligent elevator dispatch system for a multi-story building handling concurrent request queues, direction optimization, and door states.",
        skeleton: "class ElevatorController {\n  // TODO: Manage multiple elevators, requests, and dispatch algorithms\n}"
    },
    {
        id: 3,
        title: "URL Shortener",
        difficulty: "Medium",
        description: "Design a URL shortener service that encodes long URLs into short unique hashes, handles custom aliases, and tracks click analytics.",
        skeleton: "class URLService {\n  // TODO: Implement hashing and URL redirection mapping\n}\n\ninterface ShorteningStrategy {\n  String encode(long id);\n}"
    },
    {
        id: 4,
        title: "Vending Machine",
        difficulty: "Medium",
        description: "Design a state-driven vending machine managing coin insertion, product inventory, state transitions, and change return.",
        skeleton: "class VendingMachine {\n  // TODO: Manage states (Idle, HasMoney, Dispensing)\n}\n\ninterface VendingState {\n  void selectItem();\n}"
    },
    {
        id: 5,
        title: "Splitwise",
        difficulty: "Medium",
        description: "Design an expense-sharing application to track individual balances, handle equal/exact/percentage splits, and simplify debts.",
        skeleton: "class ExpenseManager {\n  // TODO: Add expenses, compute balances, simplify debts\n}\n\ninterface SplitStrategy {\n  List<Split> split(double amount, List<User> users);\n}"
    },
    {
        id: 6,
        title: "Movie Ticket Booking System",
        difficulty: "Hard",
        description: "Design a concurrent movie ticket booking platform supporting theaters, shows, seat categories, temporary locks, and payment processing.",
        skeleton: "class BookingManager {\n  // TODO: Handle concurrent seat selection, temporary locks, and payments\n}\n\ninterface PaymentStrategy {\n  boolean processPayment(double amount);\n}"
    }
];

function escapeHtml(str) {
    if (!str) return '';
    return String(str)
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#039;');
}

function updateConnectionStatus(connected, message = '') {
    const statusEl = document.getElementById('connection-status');
    const bannerEl = document.getElementById('status-banner');
    if (statusEl) {
        if (connected) {
            statusEl.innerText = 'Connected';
            statusEl.className = 'text-[10px] font-bold px-2 py-0.5 rounded-full bg-emerald-50 text-emerald-700 border border-emerald-200';
        } else {
            statusEl.innerText = 'Offline Mode';
            statusEl.className = 'text-[10px] font-bold px-2 py-0.5 rounded-full bg-amber-50 text-amber-700 border border-amber-200';
        }
    }
    if (bannerEl) {
        if (connected) {
            bannerEl.classList.add('hidden');
        } else {
            bannerEl.innerHTML = message || `⚠️ <b>Backend server unreachable at ${API_BASE || 'http://localhost:8080'}</b>.<br/>Make sure to run <code>go run .</code>. Showing local problem catalog.`;
            bannerEl.classList.remove('hidden');
        }
    }
}

async function loadProblems() {
    try {
        const res = await fetch(`${API_BASE}/api/problems`);
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        const data = await res.json();
        if (Array.isArray(data) && data.length > 0) {
            problemsData = data;
            updateConnectionStatus(true);
        } else {
            problemsData = FALLBACK_PROBLEMS;
            updateConnectionStatus(false, 'Database returned no problems. Loaded fallback catalog.');
        }
    } catch (e) {
        console.warn('Could not fetch from Go backend API, using local fallback catalog:', e);
        problemsData = FALLBACK_PROBLEMS;
        updateConnectionStatus(false);
    }

    renderProblemList();
}

function renderProblemList() {
    const listDiv = document.getElementById('problem-list');
    if (!listDiv) return;
    listDiv.innerHTML = '';

    problemsData.forEach((p) => {
        let badgeClass = 'bg-slate-100 text-slate-700 border-slate-200';
        const diff = (p.difficulty || 'Medium').toLowerCase();

        if (diff.includes('easy')) {
            badgeClass = 'bg-emerald-50 text-emerald-700 border-emerald-200';
        } else if (diff.includes('medium')) {
            badgeClass = 'bg-amber-50 text-amber-700 border-amber-200';
        } else if (diff.includes('hard')) {
            badgeClass = 'bg-rose-50 text-rose-700 border-rose-200';
        }

        const item = document.createElement('div');
        item.id = `problem-card-${p.id}`;
        item.className = 'p-3 border border-slate-200 rounded-xl cursor-pointer transition hover:border-blue-400 hover:bg-blue-50/40 bg-white group';
        item.onclick = () => selectProblem(p.id);
        item.innerHTML = `
            <div class="flex justify-between items-center">
                <span class="font-semibold text-slate-800 text-sm group-hover:text-blue-600 transition">${escapeHtml(p.title)}</span>
                <span class="text-[10px] uppercase tracking-wider px-2 py-0.5 rounded-full font-bold border ${badgeClass}">
                    ${escapeHtml(p.difficulty)}
                </span>
            </div>
        `;
        listDiv.appendChild(item);
    });

    // Select the first problem safely after all cards are in the DOM
    if (problemsData.length > 0) {
        selectProblem(problemsData[0].id);
    }
}

function selectProblem(id) {
    const p = problemsData.find(item => item.id === id);
    if (!p) return;

    selectedProblemId = p.id;

    // Highlight active card in sidebar
    problemsData.forEach(item => {
        const card = document.getElementById(`problem-card-${item.id}`);
        if (card) {
            if (item.id === id) {
                card.className = 'p-3 border-2 border-blue-500 rounded-xl cursor-pointer bg-blue-50/60 shadow-sm transition';
            } else {
                card.className = 'p-3 border border-slate-200 rounded-xl cursor-pointer transition hover:border-blue-400 hover:bg-blue-50/40 bg-white group';
            }
        }
    });

    const titleEl = document.getElementById('current-problem-title');
    if (titleEl) titleEl.innerText = p.title;

    const descEl = document.getElementById('current-problem-desc');
    if (descEl) descEl.innerText = p.description;

    const diffBadge = document.getElementById('current-problem-diff');
    if (diffBadge) {
        diffBadge.innerText = p.difficulty;
        diffBadge.className = `text-xs px-2.5 py-1 rounded-full font-bold border ${
            (p.difficulty || '').toLowerCase() === 'hard'
                ? 'bg-rose-50 text-rose-700 border-rose-200'
                : 'bg-amber-50 text-amber-700 border-amber-200'
        }`;
    }

    const codeInput = document.getElementById('code-input');
    if (codeInput) codeInput.value = p.skeleton;

    loadHistory(p.id);
}

async function submitSolution() {
    if (!selectedProblemId) return alert('Please select a problem first.');
    const code_content = document.getElementById('code-input').value;
    const design_notes = document.getElementById('design-notes').value;
    const submitBtn = document.getElementById('submit-btn');

    // UI feedback during submission
    submitBtn.disabled = true;
    submitBtn.innerHTML = `
        <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        <span>Evaluating Design...</span>
    `;

    try {
        const res = await fetch(`${API_BASE}/api/attempts`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ problem_id: selectedProblemId, code_content, design_notes })
        });

        if (!res.ok) {
            const errData = await res.json().catch(() => ({}));
            throw new Error(errData.error || `Server responded with status ${res.status}`);
        }

        const data = await res.json();
        lastResult = data.result;

        // Render feedback section
        const feedbackSec = document.getElementById('feedback-section');
        if (feedbackSec) feedbackSec.classList.remove('hidden');

        const exportBtn = document.getElementById('export-btn');
        if (exportBtn) exportBtn.classList.remove('hidden');

        // Score Badge
        const badge = document.getElementById('score-badge');
        if (badge && data.result) {
            badge.innerText = `Score: ${data.result.overall_score}/100`;
            badge.className = `px-4 py-1.5 rounded-full font-black text-base border shadow-sm ${
                data.result.overall_score >= 75
                    ? 'bg-emerald-50 text-emerald-700 border-emerald-300'
                    : data.result.overall_score >= 40
                        ? 'bg-amber-50 text-amber-700 border-amber-300'
                        : 'bg-rose-50 text-rose-700 border-rose-300'
            }`;
        }

        // Multi-dimensional Rubric Grid
        const rubricGrid = document.getElementById('rubric-grid');
        if (rubricGrid) {
            rubricGrid.innerHTML = '';
            if (data.result && Array.isArray(data.result.rubric_dimensions)) {
                data.result.rubric_dimensions.forEach(dim => {
                    rubricGrid.innerHTML += `
                        <div class="p-4 border border-slate-200 rounded-xl bg-slate-50/70 hover:bg-slate-50 transition">
                            <div class="flex justify-between items-center mb-2">
                                <h5 class="font-bold text-sm text-slate-800">${escapeHtml(dim.criterion)}</h5>
                                <span class="text-xs font-bold px-2 py-0.5 rounded ${dim.score >= 4 ? 'bg-emerald-100 text-emerald-800' : dim.score >= 2 ? 'bg-amber-100 text-amber-800' : 'bg-rose-100 text-rose-800'}">
                                    ${dim.score} / ${dim.max_score}
                                </span>
                            </div>
                            <div class="space-y-1.5 text-xs">
                                <p class="text-slate-600"><b class="text-slate-700">Evidence:</b> ${escapeHtml(dim.evidence)}</p>
                                <p class="text-amber-700"><b class="text-slate-700">Concern:</b> ${escapeHtml(dim.concern)}</p>
                                <p class="text-blue-700"><b class="text-slate-700">Suggestion:</b> ${escapeHtml(dim.suggestion)}</p>
                            </div>
                        </div>
                    `;
                });
            }
        }

        // Critique List
        const critiqueList = document.getElementById('feedback-list');
        if (critiqueList) {
            critiqueList.innerHTML = '';
            if (data.result && Array.isArray(data.result.critique)) {
                data.result.critique.forEach(item => {
                    const li = document.createElement('li');
                    li.className = 'flex items-start gap-2';
                    li.innerHTML = `<span>${escapeHtml(item)}</span>`;
                    critiqueList.appendChild(li);
                });
            }
        }

        // Render Mermaid Diagram with error boundary
        const umlDiv = document.getElementById('uml-diagram');
        if (umlDiv && data.result && data.result.mermaid_chart) {
            umlDiv.removeAttribute('data-processed');
            umlDiv.innerHTML = data.result.mermaid_chart;
            if (window.mermaid) {
                try {
                    await mermaid.run({ nodes: [umlDiv] });
                } catch (mErr) {
                    console.warn("Mermaid rendering warning:", mErr);
                }
            }
        }

        // Smooth scroll to feedback
        if (feedbackSec) feedbackSec.scrollIntoView({ behavior: 'smooth' });

        updateConnectionStatus(true);
        loadHistory(selectedProblemId);
    } catch (err) {
        console.error("Evaluation submission failed:", err);
        alert(`Evaluation failed: ${err.message || err}.\n\nEnsure the backend server is running with 'go run .' on port 8080.`);
        updateConnectionStatus(false);
    } finally {
        submitBtn.disabled = false;
        submitBtn.innerHTML = `
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path></svg>
            <span>Submit for Rubric Evaluation</span>
        `;
    }
}

function exportReport() {
    if (!lastResult) return;
    let rubricMd = "";
    if (lastResult.rubric_dimensions && Array.isArray(lastResult.rubric_dimensions)) {
        rubricMd = lastResult.rubric_dimensions.map(r => `
### ${r.criterion} (${r.score}/${r.max_score})
- **Evidence:** ${r.evidence}
- **Concern:** ${r.concern}
- **Suggestion:** ${r.suggestion}
`).join('\n');
    }

    const critiques = Array.isArray(lastResult.critique)
        ? lastResult.critique.map(c => `- ${c}`).join('\n')
        : '- Completed.';

    const report = `# Low-Level Design (LLD) Evaluation Report
- **Overall Score**: ${lastResult.overall_score}/100
- **Evaluator**: ${lastResult.evaluator_name || 'Deterministic Pipeline'}
- **Timestamp**: ${new Date().toISOString()}

## Rubric Breakdown
${rubricMd}

## Engineering Critique
${critiques}

## Architecture Diagram (Mermaid)
\`\`\`mermaid
${lastResult.mermaid_chart || ''}
\`\`\`
`;
    const blob = new Blob([report], { type: 'text/markdown' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `LLD_Evaluation_${selectedProblemId}_${Date.now()}.md`;
    a.click();
}

async function loadHistory(problemId) {
    const historyDiv = document.getElementById('history-list');
    const countSpan = document.getElementById('attempt-count');
    if (!historyDiv || !countSpan) return;

    try {
        const res = await fetch(`${API_BASE}/api/attempts/${problemId}`);
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        const attempts = await res.json();
        historyDiv.innerHTML = '';

        if (!Array.isArray(attempts) || attempts.length === 0) {
            countSpan.innerText = '0 Attempts';
            historyDiv.innerHTML = '<p class="text-slate-400 text-xs italic">No previous attempts for this problem. Submit your first design above to begin the practice loop!</p>';
            return;
        }

        countSpan.innerText = `${attempts.length} Attempt${attempts.length > 1 ? 's' : ''}`;

        attempts.forEach((a, idx) => {
            const attemptNumber = attempts.length - idx;
            const isLatest = idx === 0;
            const scoreColor = a.score >= 75
                ? 'text-emerald-600 bg-emerald-50 border-emerald-200'
                : a.score >= 40
                    ? 'text-amber-600 bg-amber-50 border-amber-200'
                    : 'text-rose-600 bg-rose-50 border-rose-200';

            const statusBadge = a.status && a.status !== 'COMPLETED'
                ? `<span class="text-[10px] bg-amber-100 text-amber-800 font-bold px-1.5 py-0.5 rounded">${escapeHtml(a.status)}</span>`
                : '';

            historyDiv.innerHTML += `
                <div class="p-3.5 border border-slate-200 rounded-xl bg-slate-50 flex flex-col sm:flex-row sm:justify-between sm:items-center gap-2 hover:bg-white transition shadow-sm">
                    <div class="space-y-1">
                        <div class="flex items-center gap-2">
                            <span class="font-bold text-xs text-slate-800">Attempt #${attemptNumber}</span>
                            ${isLatest ? '<span class="text-[10px] bg-blue-100 text-blue-700 font-bold px-1.5 py-0.5 rounded">Latest</span>' : ''}
                            ${statusBadge}
                            <span class="text-[11px] text-slate-400">${new Date(a.created_at).toLocaleString()}</span>
                        </div>
                        <p class="text-xs text-slate-500 font-mono truncate max-w-md">Notes: ${escapeHtml(a.design_notes) || 'No architectural notes provided'}</p>
                    </div>
                    <div class="flex items-center gap-3">
                        <span class="text-xs font-bold px-2.5 py-1 rounded-full border ${scoreColor}">
                            Score: ${a.score}/100
                        </span>
                        <button onclick="restoreAttempt(${a.id})" class="text-xs text-blue-600 hover:text-blue-800 font-semibold underline">
                            View / Restore
                        </button>
                    </div>
                </div>
            `;
        });

        window.currentAttempts = attempts;
    } catch (e) {
        console.warn("Could not load attempt history:", e);
        countSpan.innerText = '0 Attempts';
        historyDiv.innerHTML = '<p class="text-slate-400 text-xs italic">Attempt history is stored in the database. Run the backend server with <code>go run .</code> to persist and view attempts.</p>';
    }
}

function restoreAttempt(attemptId) {
    if (!window.currentAttempts) return;
    const attempt = window.currentAttempts.find(a => a.id === attemptId);
    if (!attempt) return;

    const codeInput = document.getElementById('code-input');
    if (codeInput) codeInput.value = attempt.code_content || '';

    const notesInput = document.getElementById('design-notes');
    if (notesInput) notesInput.value = attempt.design_notes || '';

    if (attempt.result) {
        lastResult = attempt.result;
        const feedbackSec = document.getElementById('feedback-section');
        if (feedbackSec) feedbackSec.classList.remove('hidden');

        const exportBtn = document.getElementById('export-btn');
        if (exportBtn) exportBtn.classList.remove('hidden');

        const badge = document.getElementById('score-badge');
        if (badge) {
            badge.innerText = `Score: ${attempt.result.overall_score}/100`;
            badge.className = `px-4 py-1.5 rounded-full font-black text-base border shadow-sm ${
                attempt.result.overall_score >= 75 ? 'bg-emerald-50 text-emerald-700 border-emerald-300' : 'bg-amber-50 text-amber-700 border-amber-300'
            }`;
        }

        const rubricGrid = document.getElementById('rubric-grid');
        if (rubricGrid && Array.isArray(attempt.result.rubric_dimensions)) {
            rubricGrid.innerHTML = '';
            attempt.result.rubric_dimensions.forEach(dim => {
                rubricGrid.innerHTML += `
                    <div class="p-4 border border-slate-200 rounded-xl bg-slate-50/70">
                        <div class="flex justify-between items-center mb-2">
                            <h5 class="font-bold text-sm text-slate-800">${escapeHtml(dim.criterion)}</h5>
                            <span class="text-xs font-bold px-2 py-0.5 rounded ${dim.score >= 4 ? 'bg-emerald-100 text-emerald-800' : dim.score >= 2 ? 'bg-amber-100 text-amber-800' : 'bg-rose-100 text-rose-800'}">
                                ${dim.score} / ${dim.max_score}
                            </span>
                        </div>
                        <div class="space-y-1.5 text-xs">
                            <p class="text-slate-600"><b class="text-slate-700">Evidence:</b> ${escapeHtml(dim.evidence)}</p>
                            <p class="text-amber-700"><b class="text-slate-700">Concern:</b> ${escapeHtml(dim.concern)}</p>
                            <p class="text-blue-700"><b class="text-slate-700">Suggestion:</b> ${escapeHtml(dim.suggestion)}</p>
                        </div>
                    </div>
                `;
            });
        }

        const umlDiv = document.getElementById('uml-diagram');
        if (umlDiv && attempt.result.mermaid_chart) {
            umlDiv.removeAttribute('data-processed');
            umlDiv.innerHTML = attempt.result.mermaid_chart;
            if (window.mermaid) {
                try {
                    mermaid.run({ nodes: [umlDiv] });
                } catch (mErr) {
                    console.warn("Mermaid restore warning:", mErr);
                }
            }
        }

        if (feedbackSec) feedbackSec.scrollIntoView({ behavior: 'smooth' });
    }
}

// Automatically load problems on startup
loadProblems();