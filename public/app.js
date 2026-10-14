let selectedProblemId = null;
let problemsData = [];
let lastResult = null;

async function loadProblems() {
    try {
        const res = await fetch('/api/problems');
        problemsData = await res.json();
        const listDiv = document.getElementById('problem-list');
        listDiv.innerHTML = '';
        
        problemsData.forEach((p, idx) => {
            let badgeClass = 'bg-slate-100 text-slate-700 border-slate-200';
            const diff = p.difficulty.toLowerCase();
            
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
                    <span class="font-semibold text-slate-800 text-sm group-hover:text-blue-600 transition">${p.title}</span>
                    <span class="text-[10px] uppercase tracking-wider px-2 py-0.5 rounded-full font-bold border ${badgeClass}">
                        ${p.difficulty}
                    </span>
                </div>
            `;
            listDiv.appendChild(item);
            if (idx === 0) selectProblem(p.id);
        });
    } catch (e) {
        console.error("Failed to load problems", e);
    }
}

function selectProblem(id) {
    const p = problemsData.find(item => item.id === id);
    if (!p) return;

    selectedProblemId = p.id;
    
    // Highlight active in list
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

    document.getElementById('current-problem-title').innerText = p.title;
    document.getElementById('current-problem-desc').innerText = p.description;
    
    const diffBadge = document.getElementById('current-problem-diff');
    diffBadge.innerText = p.difficulty;
    diffBadge.className = `text-xs px-2.5 py-1 rounded-full font-bold border ${
        p.difficulty.toLowerCase() === 'hard' ? 'bg-rose-50 text-rose-700 border-rose-200' : 'bg-amber-50 text-amber-700 border-amber-200'
    }`;

    document.getElementById('code-input').value = p.skeleton;
    loadHistory(p.id);
}

async function submitSolution() {
    if (!selectedProblemId) return alert('Please select a problem first.');
    const code_content = document.getElementById('code-input').value;
    const design_notes = document.getElementById('design-notes').value;
    const submitBtn = document.getElementById('submit-btn');

    // UI feedback for submission
    submitBtn.disabled = true;
    submitBtn.innerHTML = `
        <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        <span>Evaluating Design...</span>
    `;

    try {
        const res = await fetch('/api/attempts', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ problem_id: selectedProblemId, code_content, design_notes })
        });
        const data = await res.json();
        lastResult = data.result;
        
        // Render feedback section
        const feedbackSec = document.getElementById('feedback-section');
        feedbackSec.classList.remove('hidden');
        document.getElementById('export-btn').classList.remove('hidden');

        // Score Badge
        const badge = document.getElementById('score-badge');
        badge.innerText = `Score: ${data.result.overall_score}/100`;
        badge.className = `px-4 py-1.5 rounded-full font-black text-base border shadow-sm ${
            data.result.overall_score >= 75 
                ? 'bg-emerald-50 text-emerald-700 border-emerald-300' 
                : data.result.overall_score >= 40 
                    ? 'bg-amber-50 text-amber-700 border-amber-300' 
                    : 'bg-rose-50 text-rose-700 border-rose-300'
        }`;

        // Populate Rubric Grid
        const rubricGrid = document.getElementById('rubric-grid');
        rubricGrid.innerHTML = '';
        if (data.result.rubric_dimensions && data.result.rubric_dimensions.length > 0) {
            data.result.rubric_dimensions.forEach(dim => {
                rubricGrid.innerHTML += `
                    <div class="p-4 border border-slate-200 rounded-xl bg-slate-50/70 hover:bg-slate-50 transition">
                        <div class="flex justify-between items-center mb-2">
                            <h5 class="font-bold text-sm text-slate-800">${dim.criterion}</h5>
                            <span class="text-xs font-bold px-2 py-0.5 rounded ${dim.score >= 4 ? 'bg-emerald-100 text-emerald-800' : dim.score >= 2 ? 'bg-amber-100 text-amber-800' : 'bg-rose-100 text-rose-800'}">
                                ${dim.score} / ${dim.max_score}
                            </span>
                        </div>
                        <div class="space-y-1.5 text-xs">
                            <p class="text-slate-600"><b class="text-slate-700">Evidence:</b> ${dim.evidence}</p>
                            <p class="text-amber-700"><b class="text-slate-700">Concern:</b> ${dim.concern}</p>
                            <p class="text-blue-700"><b class="text-slate-700">Suggestion:</b> ${dim.suggestion}</p>
                        </div>
                    </div>
                `;
            });
        }

        // Populate critique list
        const critiqueList = document.getElementById('feedback-list');
        critiqueList.innerHTML = '';
        data.result.critique.forEach(item => {
            const li = document.createElement('li');
            li.className = 'flex items-start gap-2';
            li.innerHTML = `<span>${item}</span>`;
            critiqueList.appendChild(li);
        });

        // Render Mermaid diagram
        const umlDiv = document.getElementById('uml-diagram');
        umlDiv.removeAttribute('data-processed');
        umlDiv.innerHTML = data.result.mermaid_chart;
        mermaid.run({ nodes: [umlDiv] });

        // Scroll to feedback smoothly
        feedbackSec.scrollIntoView({ behavior: 'smooth' });

        loadHistory(selectedProblemId);
    } catch (err) {
        alert("Failed to submit design: " + err);
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
    if (lastResult.rubric_dimensions) {
        rubricMd = lastResult.rubric_dimensions.map(r => `
### ${r.criterion} (${r.score}/${r.max_score})
- **Evidence:** ${r.evidence}
- **Concern:** ${r.concern}
- **Suggestion:** ${r.suggestion}
`).join('\n');
    }

    const report = `# Low-Level Design (LLD) Evaluation Report
- **Overall Score**: ${lastResult.overall_score}/100
- **Evaluator**: ${lastResult.evaluator_name || 'Deterministic Pipeline'}
- **Timestamp**: ${new Date().toISOString()}

## Rubric Breakdown
${rubricMd}

## Engineering Critique
${lastResult.critique.map(c => `- ${c}`).join('\n')}

## Architecture Diagram (Mermaid)
\`\`\`mermaid
${lastResult.mermaid_chart}
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
    try {
        const res = await fetch(`/api/attempts/${problemId}`);
        const attempts = await res.json();
        const historyDiv = document.getElementById('history-list');
        const countSpan = document.getElementById('attempt-count');
        historyDiv.innerHTML = '';

        if (!attempts || attempts.length === 0) {
            countSpan.innerText = '0 Attempts';
            historyDiv.innerHTML = '<p class="text-slate-400 text-xs italic">No previous attempts for this problem. Submit your first design above to begin the practice loop!</p>';
            return;
        }

        countSpan.innerText = `${attempts.length} Attempt${attempts.length > 1 ? 's' : ''}`;

        attempts.forEach((a, idx) => {
            const attemptNumber = attempts.length - idx;
            const isLatest = idx === 0;
            const scoreColor = a.score >= 75 ? 'text-emerald-600 bg-emerald-50 border-emerald-200' : a.score >= 40 ? 'text-amber-600 bg-amber-50 border-amber-200' : 'text-rose-600 bg-rose-50 border-rose-200';
            
            historyDiv.innerHTML += `
                <div class="p-3.5 border border-slate-200 rounded-xl bg-slate-50 flex flex-col sm:flex-row sm:justify-between sm:items-center gap-2 hover:bg-white transition shadow-sm">
                    <div class="space-y-1">
                        <div class="flex items-center gap-2">
                            <span class="font-bold text-xs text-slate-800">Attempt #${attemptNumber}</span>
                            ${isLatest ? '<span class="text-[10px] bg-blue-100 text-blue-700 font-bold px-1.5 py-0.5 rounded">Latest</span>' : ''}
                            <span class="text-[11px] text-slate-400">${new Date(a.created_at).toLocaleString()}</span>
                        </div>
                        <p class="text-xs text-slate-500 font-mono truncate max-w-md">Notes: ${a.design_notes || 'No architectural notes provided'}</p>
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

        // Store globally for quick restore
        window.currentAttempts = attempts;
    } catch (e) {
        console.error("Failed to load history", e);
    }
}

function restoreAttempt(attemptId) {
    if (!window.currentAttempts) return;
    const attempt = window.currentAttempts.find(a => a.id === attemptId);
    if (!attempt) return;

    document.getElementById('code-input').value = attempt.code_content;
    document.getElementById('design-notes').value = attempt.design_notes;

    if (attempt.result) {
        lastResult = attempt.result;
        const feedbackSec = document.getElementById('feedback-section');
        feedbackSec.classList.remove('hidden');
        document.getElementById('export-btn').classList.remove('hidden');

        const badge = document.getElementById('score-badge');
        badge.innerText = `Score: ${attempt.result.overall_score}/100`;
        badge.className = `px-4 py-1.5 rounded-full font-black text-base border shadow-sm ${
            attempt.result.overall_score >= 75 ? 'bg-emerald-50 text-emerald-700 border-emerald-300' : 'bg-amber-50 text-amber-700 border-amber-300'
        }`;

        // Populate Rubric Grid
        const rubricGrid = document.getElementById('rubric-grid');
        rubricGrid.innerHTML = '';
        if (attempt.result.rubric_dimensions) {
            attempt.result.rubric_dimensions.forEach(dim => {
                rubricGrid.innerHTML += `
                    <div class="p-4 border border-slate-200 rounded-xl bg-slate-50/70">
                        <div class="flex justify-between items-center mb-2">
                            <h5 class="font-bold text-sm text-slate-800">${dim.criterion}</h5>
                            <span class="text-xs font-bold px-2 py-0.5 rounded ${dim.score >= 4 ? 'bg-emerald-100 text-emerald-800' : 'bg-amber-100 text-amber-800'}">
                                ${dim.score} / ${dim.max_score}
                            </span>
                        </div>
                        <div class="space-y-1.5 text-xs">
                            <p class="text-slate-600"><b class="text-slate-700">Evidence:</b> ${dim.evidence}</p>
                            <p class="text-amber-700"><b class="text-slate-700">Concern:</b> ${dim.concern}</p>
                            <p class="text-blue-700"><b class="text-slate-700">Suggestion:</b> ${dim.suggestion}</p>
                        </div>
                    </div>
                `;
            });
        }

        // Mermaid
        const umlDiv = document.getElementById('uml-diagram');
        umlDiv.removeAttribute('data-processed');
        umlDiv.innerHTML = attempt.result.mermaid_chart;
        mermaid.run({ nodes: [umlDiv] });

        feedbackSec.scrollIntoView({ behavior: 'smooth' });
    }
}

loadProblems();