import { GetBestTimes } from '../wailsjs/go/main/App';

export interface DifficultyOption {
  value: 'easy' | 'medium' | 'hard';
  label: string;
  boardSize: number;
  totalCells: number;
  mines: number;
}

export const DIFFICULTIES: DifficultyOption[] = [
  { value: 'easy', label: 'Fácil', boardSize: 5, totalCells: 150, mines: 15 },
  { value: 'medium', label: 'Médio', boardSize: 7, totalCells: 294, mines: 40 },
  { value: 'hard', label: 'Difícil', boardSize: 9, totalCells: 486, mines: 80 },
];

const DEFAULT_DIFFICULTY = 'medium';

export function formatBestTime(seconds: number | undefined | null): string {
  if (seconds === undefined || seconds === null) {
    return '—';
  }
  const minutes = Math.floor(seconds / 60);
  const rest = seconds % 60;
  return `${String(minutes).padStart(2, '0')}:${String(rest).padStart(2, '0')}`;
}

function minePositions(option: DifficultyOption): Set<number> {
  const cells = option.boardSize * option.boardSize;
  const count = Math.round((option.mines / option.totalCells) * cells);
  const positions = new Set<number>();
  let seed = option.boardSize * 7919;
  while (positions.size < count) {
    seed = (seed * 1103515245 + 12345) % 2147483648;
    positions.add(seed % cells);
  }
  return positions;
}

function renderMinefield(option: DifficultyOption): string {
  const cells = option.boardSize * option.boardSize;
  const mines = minePositions(option);
  const squares = Array.from({ length: cells }, (_, i) =>
    `<i class="${mines.has(i) ? 'cell cell-mine' : 'cell'}"></i>`,
  ).join('');
  return `<div class="minefield" style="--n:${option.boardSize}" aria-hidden="true">${squares}</div>`;
}

function renderCard(option: DifficultyOption): string {
  const density = Math.round((option.mines / option.totalCells) * 100);
  return `
    <button class="card" type="button" role="radio" aria-checked="false"
            id="difficulty-${option.value}" data-difficulty="${option.value}">
      ${renderMinefield(option)}
      <span class="card-label">${option.label}</span>
      <span class="card-specs">
        <span>${option.boardSize}×${option.boardSize} por face</span>
        <span>${option.totalCells} células</span>
        <span>${option.mines} minas · ${density}%</span>
      </span>
      <span class="card-record">
        <span class="card-record-label">recorde</span>
        <span class="card-record-value" data-record="${option.value}">—</span>
      </span>
    </button>
  `;
}

export interface MenuHandle {
  hide(): void;
}

export function mountMenu(
  container: HTMLElement,
  onPlay: (option: DifficultyOption) => Promise<void>,
): MenuHandle {
  container.innerHTML = `
    <div class="menu-inner">
      <header class="menu-header">
        <h1 class="menu-title">Cubo Minado</h1>
        <p class="menu-tagline">Campo minado nas seis faces de um cubo.</p>
      </header>
      <div class="cards" role="radiogroup" aria-label="Dificuldade">
        ${DIFFICULTIES.map(renderCard).join('')}
      </div>
      <div class="menu-actions">
        <button class="play" type="button" id="play">Jogar</button>
        <p class="menu-error" id="menu-error" role="alert" hidden></p>
      </div>
    </div>
  `;

  const cards = Array.from(container.querySelectorAll<HTMLButtonElement>('.card'));
  const playButton = container.querySelector('#play') as HTMLButtonElement;
  const errorMessage = container.querySelector('#menu-error') as HTMLParagraphElement;

  let selected =
    DIFFICULTIES.find((option) => option.value === DEFAULT_DIFFICULTY) ?? DIFFICULTIES[0];

  const applySelection = () => {
    for (const card of cards) {
      const isSelected = card.dataset.difficulty === selected.value;
      card.classList.toggle('selected', isSelected);
      card.setAttribute('aria-checked', String(isSelected));
      card.tabIndex = isSelected ? 0 : -1;
    }
  };

  for (const card of cards) {
    card.addEventListener('click', () => {
      const option = DIFFICULTIES.find((d) => d.value === card.dataset.difficulty);
      if (option) {
        selected = option;
        errorMessage.hidden = true;
        applySelection();
      }
    });
  }

  applySelection();

  playButton.addEventListener('click', () => {
    errorMessage.hidden = true;
    playButton.disabled = true;
    void onPlay(selected).catch((err) => {
      errorMessage.textContent = `Não foi possível iniciar a partida: ${
        err instanceof Error ? err.message : String(err)
      }`;
      errorMessage.hidden = false;
    }).finally(() => {
      playButton.disabled = false;
    });
  });

  void loadBestTimes(container);

  return {
    hide() {
      container.hidden = true;
    },
  };
}

async function loadBestTimes(container: HTMLElement): Promise<void> {
  let times;
  try {
    times = await GetBestTimes();
  } catch {
    return;
  }

  const byDifficulty: Record<string, number | undefined> = {
    easy: times.easy,
    medium: times.medium,
    hard: times.hard,
  };

  for (const option of DIFFICULTIES) {
    const slot = container.querySelector(`[data-record="${option.value}"]`);
    if (slot) {
      slot.textContent = formatBestTime(byDifficulty[option.value]);
    }
  }
}
