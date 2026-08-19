import {
  ChordCell,
  FlagCell,
  GetBestTimes,
  NewGame,
  OpenCell,
} from '../wailsjs/go/main/App';

const DIFFICULTIES = ['easy', 'medium', 'hard'] as const;

export function mountDebugPanel(container: HTMLElement): void {
  container.innerHTML = `
    <h2>Debug — bindings (temporário)</h2>
    <div class="buttons">
      <select id="debug-difficulty">
        ${DIFFICULTIES.map((d) => `<option value="${d}">${d}</option>`).join('')}
        <option value="impossible">impossible (inválida)</option>
      </select>
      <button id="debug-new-game" type="button">NewGame</button>
      <button id="debug-open-cell" type="button">OpenCell</button>
      <button id="debug-flag-cell" type="button">FlagCell</button>
      <button id="debug-chord-cell" type="button">ChordCell</button>
      <button id="debug-best-times" type="button">GetBestTimes</button>
    </div>
    <pre id="debug-output">Nenhuma chamada ainda.</pre>
  `;

  const output = container.querySelector('#debug-output') as HTMLPreElement;
  const difficulty = container.querySelector('#debug-difficulty') as HTMLSelectElement;

  const invoke = async (label: string, call: () => Promise<unknown>) => {
    output.classList.remove('error');
    output.textContent = `${label}: chamando...`;
    try {
      const result = await call();
      output.textContent = `${label} →\n${JSON.stringify(result, null, 2)}`;
    } catch (err) {
      output.classList.add('error');
      output.textContent = `${label} ✗ ${err instanceof Error ? err.message : String(err)}`;
    }
  };

  const bind = (id: string, label: string, call: () => Promise<unknown>) => {
    (container.querySelector(`#${id}`) as HTMLButtonElement).addEventListener('click', () => {
      void invoke(label, call);
    });
  };

  bind('debug-new-game', 'NewGame', () => NewGame(difficulty.value));
  bind('debug-open-cell', 'OpenCell', () => OpenCell(0, 0, 0));
  bind('debug-flag-cell', 'FlagCell', () => FlagCell(0, 0, 0));
  bind('debug-chord-cell', 'ChordCell', () => ChordCell(0, 0, 0));
  bind('debug-best-times', 'GetBestTimes', () => GetBestTimes());
}
