import { NewGame } from '../wailsjs/go/main/App';
import { mountMenu, type DifficultyOption } from './menu';
import { createPlaceholderScene, type PlaceholderScene } from './three/scene';

export interface AppShellContainers {
  menu: HTMLElement;
  scene: HTMLElement;
  status: HTMLElement;
}

export function mountAppShell(containers: AppShellContainers): void {
  let scene: PlaceholderScene | null = null;

  const startMatch = async (option: DifficultyOption) => {
    const board = await NewGame(option.value);

    menu.hide();
    containers.scene.hidden = false;

    if (!scene) {
      scene = createPlaceholderScene(containers.scene);
    }
    scene.start();

    containers.status.innerHTML = `
      <span class="status-label">Tabuleiro gerado</span>
      <span class="status-value">${option.label} · ${board.boardSize}×${board.boardSize} por face</span>
    `;
    containers.status.hidden = false;
  };

  const menu = mountMenu(containers.menu, startMatch);

  containers.scene.hidden = true;
  containers.status.hidden = true;
}
