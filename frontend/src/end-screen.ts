import { formatBestTime } from './menu';

export interface EndScreenActions {
  onPlayAgain(): void;
  onMainMenu(): void;
}

export interface VictoryPayload {
  result: 'victory';
  finalSeconds: number;
  previousBest: number | null;
  isNewRecord: boolean;
}

export interface DefeatPayload {
  result: 'defeat';
}

export type EndScreenPayload = VictoryPayload | DefeatPayload;

export interface EndScreen {
  show(payload: EndScreenPayload): void;
  hide(): void;
}

function victoryBody(payload: VictoryPayload): string {
  const record = payload.isNewRecord
    ? '<p class="end-record">Novo recorde</p>'
    : '';

  return `
    <p class="end-eyebrow">Vitória</p>
    <h2 class="end-title">Cubo limpo</h2>
    <dl class="end-stats">
      <div class="end-stat">
        <dt>Seu tempo</dt>
        <dd>${formatBestTime(payload.finalSeconds)}</dd>
      </div>
      <div class="end-stat">
        <dt>Recorde anterior</dt>
        <dd>${formatBestTime(payload.previousBest)}</dd>
      </div>
    </dl>
    ${record}
  `;
}

function defeatBody(): string {
  return `
    <p class="end-eyebrow">Derrota</p>
    <h2 class="end-title">Você abriu uma mina</h2>
    <p class="end-hint">O cubo abaixo mostra onde estavam todas as minas. Gire para conferir.</p>
  `;
}

export function mountEndScreen(container: HTMLElement, actions: EndScreenActions): EndScreen {
  container.hidden = true;

  const render = (payload: EndScreenPayload) => {
    container.innerHTML = `
      <div class="end-card end-card-${payload.result}">
        ${payload.result === 'victory' ? victoryBody(payload) : defeatBody()}
        <div class="end-actions">
          <button class="end-button end-button-primary" type="button" id="end-play-again">
            Jogar novamente
          </button>
          <button class="end-button" type="button" id="end-main-menu">Menu principal</button>
        </div>
      </div>
    `;

    (container.querySelector('#end-play-again') as HTMLButtonElement).addEventListener(
      'click',
      () => actions.onPlayAgain(),
    );
    (container.querySelector('#end-main-menu') as HTMLButtonElement).addEventListener(
      'click',
      () => actions.onMainMenu(),
    );
  };

  return {
    show(payload) {
      render(payload);
      container.hidden = false;
      (container.querySelector('#end-play-again') as HTMLButtonElement).focus();
    },
    hide() {
      container.hidden = true;
      container.innerHTML = '';
    },
  };
}
