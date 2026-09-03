import { DIFFICULTIES } from './menu';
import type { RenderBoard } from './three/cell-materials';

export interface HudActions {
  onRestart(): void;
  onMenu(): void;
}

export interface Hud {
  update(board: RenderBoard): void;
  setActionsHidden(hidden: boolean): void;
  dispose(): void;
}

function formatClock(seconds: number): string {
  const safe = Math.max(0, seconds);
  const minutes = Math.floor(safe / 60);
  const rest = safe % 60;
  return `${String(minutes).padStart(2, '0')}:${String(rest).padStart(2, '0')}`;
}

function totalMinesFor(board: RenderBoard): number {
  return DIFFICULTIES.find((option) => option.value === board.difficulty)?.mines ?? 0;
}

function countFlags(board: RenderBoard): number {
  let flags = 0;
  for (const face of board.faces) {
    for (const row of face.cells) {
      for (const cell of row) {
        if (cell.state === 'flagged') {
          flags++;
        }
      }
    }
  }
  return flags;
}

export function elapsedSeconds(board: RenderBoard): number {
  if (board.startedAt === undefined || board.startedAt === null) {
    return 0;
  }
  if (board.endedAt !== undefined && board.endedAt !== null) {
    return board.endedAt - board.startedAt;
  }
  return Math.floor(Date.now() / 1000) - board.startedAt;
}

function isRunning(board: RenderBoard): boolean {
  return (
    board.startedAt !== undefined &&
    board.startedAt !== null &&
    (board.endedAt === undefined || board.endedAt === null)
  );
}

export function mountHud(container: HTMLElement, actions: HudActions): Hud {
  container.innerHTML = `
    <div class="hud-stat">
      <span class="hud-label">Minas</span>
      <span class="hud-value" id="hud-mines">—</span>
    </div>
    <div class="hud-stat">
      <span class="hud-label">Tempo</span>
      <span class="hud-value" id="hud-timer">00:00</span>
    </div>
    <div class="hud-actions">
      <button class="hud-button" type="button" id="hud-restart">Reiniciar</button>
      <button class="hud-button" type="button" id="hud-menu">Menu</button>
    </div>
  `;

  const minesSlot = container.querySelector('#hud-mines') as HTMLSpanElement;
  const timerSlot = container.querySelector('#hud-timer') as HTMLSpanElement;
  const actionsRow = container.querySelector('.hud-actions') as HTMLDivElement;

  (container.querySelector('#hud-restart') as HTMLButtonElement).addEventListener(
    'click',
    () => actions.onRestart(),
  );
  (container.querySelector('#hud-menu') as HTMLButtonElement).addEventListener(
    'click',
    () => actions.onMenu(),
  );

  let current: RenderBoard | null = null;
  let ticker = 0;

  const renderTimer = () => {
    if (current) {
      timerSlot.textContent = formatClock(elapsedSeconds(current));
    }
  };

  const stopTicking = () => {
    if (ticker !== 0) {
      window.clearInterval(ticker);
      ticker = 0;
    }
  };

  return {
    update(board) {
      current = board;
      minesSlot.textContent = String(totalMinesFor(board) - countFlags(board));
      renderTimer();

      if (isRunning(board)) {
        if (ticker === 0) {
          ticker = window.setInterval(renderTimer, 250);
        }
      } else {
        stopTicking();
      }
    },
    setActionsHidden(hidden) {
      actionsRow.hidden = hidden;
    },
    dispose() {
      stopTicking();
      current = null;
      container.innerHTML = '';
    },
  };
}
