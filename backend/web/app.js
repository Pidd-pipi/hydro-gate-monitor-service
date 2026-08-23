const statusNode = document.querySelector('#status');
const list = document.querySelector('#gates');
async function loadGates() {
  const response = await fetch('/api/gates');
  const data = await response.json();
  list.replaceChildren(...data.gates.map(gate => {
    const item = document.createElement('li');
    item.textContent = `${gate.name} - ${gate.state} - alert: ${gate.alertLevel}`;
    if (!gate.acknowledged && gate.alertLevel !== 'normal') {
      const button = document.createElement('button');
      button.textContent = 'Acknowledge';
      button.onclick = async () => { await fetch(`/api/gates/${gate.id}/ack`, {method:'POST', headers:{'Content-Type':'application/json'}, body:JSON.stringify({note:'Reviewed in web console'})}); loadGates(); };
      item.append(' ', button);
    }
    return item;
  }));
  statusNode.textContent = `${data.gates.length} gates observed`;
}
loadGates().catch(error => { statusNode.textContent = error.message; });
