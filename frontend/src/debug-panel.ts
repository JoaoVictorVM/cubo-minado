import {
  ChordCell,
  FlagCell,
  GetBestTimes,
  NewGame,
  OpenCell,
  SubmitTime,
} from '../wailsjs/go/main/App';
import { buildCellStatePreviewBoard } from './three/preview-board';
import type { RenderBoard } from './three/cell-materials';

const DIFFICULTIES = ['easy', 'medium', 'hard'] as const;

export interface DebugPanelHooks {
  showBoard(board: RenderBoard): void;
}

export function mountDebugPanel(container: HTMLElement, hooks: DebugPanelHooks): void {
  container.innerHTML = `
    <button class="debug-toggle" type="button" id="debug-toggle" aria-expanded="false">
      Debug — bindings (temporário)
    </button>
    <div class="debug-body">
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
    <div class="buttons">
      <input id="debug-seconds" type="number" min="0" step="1" value="60" />
      <button id="debug-submit-time" type="button">SubmitTime</button>
    </div>
    <div class="buttons">
      <button id="debug-preview-cells" type="button">Prever estados das células</button>
    </div>
    <pre id="debug-output">Nenhuma chamada ainda.</pre>
    </div>
  `;

  container.classList.add('collapsed');
  const toggle = container.querySelector('#debug-toggle') as HTMLButtonElement;
  toggle.addEventListener('click', () => {
    const collapsed = container.classList.toggle('collapsed');
    toggle.setAttribute('aria-expanded', String(!collapsed));
  });

  const output = container.querySelector('#debug-output') as HTMLPreElement;
  const difficulty = container.querySelector('#debug-difficulty') as HTMLSelectElement;
  const seconds = container.querySelector('#debug-seconds') as HTMLInputElement;

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
  bind('debug-submit-time', 'SubmitTime', () =>
    SubmitTime(difficulty.value, Number.parseInt(seconds.value, 10)),
  );

  const previewButton = container.querySelector('#debug-preview-cells') as HTMLButtonElement;
  previewButton.addEventListener('click', () => {
    const board = buildCellStatePreviewBoard();
    hooks.showBoard(board);
    output.classList.remove('error');
    output.textContent = 'Prévia de estados renderizada na face frontal do cubo.';
  });
}
