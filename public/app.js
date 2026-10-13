let selectedProblemId = null;
let problemsData = [];
let lastResult = null;

async function loadProblems() {
    const res = await fetch('/api/problems');
    problemsData = await res.json();
    const listDiv = document.getElementById('problem-list');
    listDiv.innerHTML = '';
    
    problemsData.forEach((p, idx) => {
        let badgeColor = 'bg-gray-100 text-gray-800';
        let hoverColor = '';
        const diff = p.difficulty.toLowerCase();
        
        if (diff.includes('easy')) {
            badgeColor = 'bg-green-100 text-green-800 border-green-300';
            hoverColor = 'hover:bg-green-200 hover:border-green-500';
        } else if (diff.includes('medium')) {
            badgeColor = 'bg-yellow-100 text-yellow-800 border-yellow-300';
            hoverColor = 'hover:bg-yellow-200 hover:border-yellow-500';
        } else if (diff.includes('hard')) {
            badgeColor = 'bg-red-100 text-red-800 border-red-300';
            hoverColor = 'hover:bg-red-200 hover:border-red-500';
        }

        listDiv.innerHTML += `
            <div onclick="selectProblem(${p.id})" class="p-3 border rounded-lg cursor-pointer transition ${hoverColor} group relative">
                <div class="flex justify-between items-center">
                    <span class="font-semibold text-gray-800">${p.title}</span>
                    <span class="text-xs px-2.5 py-1 rounded-full font-bold border transition ${badgeColor}" title="Difficulty: ${p.difficulty}">
                        ${p.difficulty}
                    </span>
                </div>
            </div>
        `;
        if (idx === 0) selectProblem(p.id);
    });
}

function selectProblem(id) {
    const p = problemsData.find(item => item.id === id);
    if (!p) return;

    selectedProblemId = p.id;
    document.getElementById('current-problem-title').innerText = p.title;
    document.getElementById('current-problem-desc').innerText = p.description;
    document.getElementById('code-input').value = p.skeleton;
    loadHistory(p.id);
}

async function submitSolution() {
    if (!selectedProblemId) return alert('Please select a problem first.');
    const code_content = document.getElementById('code-input').value;
    const design_notes = document.getElementById('design-notes').value;

    const res = await fetch('/api/attempts', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ problem_id: selectedProblemId, code_content, design_notes })
    });
    const data = await res.json();
    lastResult = data.result;
    
    // Show Feedback Section
    const feedbackSec = document.getElementById('feedback-section');
    feedbackSec.classList.remove('hidden');
    document.getElementById('export-btn').classList.remove('hidden');

    const badge = document.getElementById('score-badge');
    badge.innerText = `Overall Score: ${data.result.overall_score}/100`;
    badge.className = `px-4 py-1.5 rounded-full font-bold text-base ${data.result.overall_score >= 80 ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'}`;

    // Update Progress bars
    document.getElementById('bar-solid').style.width = data.result.solid_score + '%';
    document.getElementById('val-solid').innerText = data.result.solid_score + '%';

    document.getElementById('bar-pattern').style.width = data.result.pattern_score + '%';
    document.getElementById('val-pattern').innerText = data.result.pattern_score + '%';

    document.getElementById('bar-concurrency').style.width = data.result.concurrency_score + '%';
    document.getElementById('val-concurrency').innerText = data.result.concurrency_score + '%';
    
    const critiqueList = document.getElementById('feedback-list');
    critiqueList.innerHTML = '';
    data.result.critique.forEach(item => {
        critiqueList.innerHTML += `<li>${item}</li>`;
    });

    // Render Mermaid diagram
    const umlDiv = document.getElementById('uml-diagram');
    umlDiv.removeAttribute('data-processed');
    umlDiv.innerHTML = data.result.mermaid_chart;
    mermaid.run({ nodes: [umlDiv] });

    loadHistory(selectedProblemId);
}

function exportReport() {
    if (!lastResult) return;
    const report = `# LLD Evaluation Report\n\n- **Overall Score**: ${lastResult.overall_score}/100\n- **SOLID Score**: ${lastResult.solid_score}%\n- **Design Patterns**: ${lastResult.pattern_score}%\n- **Concurrency Safety**: ${lastResult.concurrency_score}%\n\n### Critique:\n${lastResult.critique.map(c => `- ${c}`).join('\n')}`;
    const blob = new Blob([report], { type: 'text/markdown' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'LLD_Evaluation_Report.md';
    a.click();
}

async function loadHistory(problemId) {
    const res = await fetch(`/api/attempts/${problemId}`);
    const attempts = await res.json();
    const historyDiv = document.getElementById('history-list');
    historyDiv.innerHTML = '';
    if (!attempts || attempts.length === 0) {
        historyDiv.innerHTML = '<p class="text-gray-500">No previous attempts yet.</p>';
        return;
    }
    attempts.forEach(a => {
        historyDiv.innerHTML += `
            <div class="p-2.5 border rounded-lg flex justify-between items-center bg-gray-50">
                <span>Score: <b>${a.score}/100</b> - ${new Date(a.created_at).toLocaleString()}</span>
                <span class="text-xs text-gray-500">Notes: ${a.design_notes.substring(0, 35)}...</span>
            </div>
        `;
    });
}

loadProblems();