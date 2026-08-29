import { NewGame } from '../wailsjs/go/main/App';
import { mountMenu, type DifficultyOption } from './menu';
import type { RenderBoard } from './three/cell-materials';
import { createCubeScene, type CubeScene } from './three/scene';

export interface AppShellContainers {
  menu: HTMLElement;
  scene: HTMLElement;
  status: HTMLElement;
}

export interface AppShellHandle {
  showBoard(board: RenderBoard, label: string): void;
}

export function mountAppShell(containers: AppShellContainers): AppShellHandle {
  let scene: CubeScene | null = null;

  const showBoard = (board: RenderBoard, label: string) => {
    menu.hide();
    containers.scene.hidden = false;

    if (scene) {
      scene.setBoard(board);
    } else {
      scene = createCubeScene(containers.scene, board);
    }
    scene.start();

    containers.status.innerHTML = `
      <span class="status-label">Tabuleiro gerado</span>
      <span class="status-value">${label} · ${board.boardSize}×${board.boardSize} por face</span>
    `;
    containers.status.hidden = false;
  };

  const startMatch = async (option: DifficultyOption) => {
    const board = await NewGame(option.value);
    showBoard(board, option.label);
  };

  const menu = mountMenu(containers.menu, startMatch);

  containers.scene.hidden = true;
  containers.status.hidden = true;

  return { showBoard };
}
