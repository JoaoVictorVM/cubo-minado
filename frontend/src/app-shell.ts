import { NewGame } from '../wailsjs/go/main/App';
import { mountHud, type Hud } from './hud';
import { mountMenu, type DifficultyOption } from './menu';
import type { RenderBoard } from './three/cell-materials';
import { createCubeScene, type CubeScene } from './three/scene';

export interface AppShellContainers {
  menu: HTMLElement;
  scene: HTMLElement;
  status: HTMLElement;
}

export interface AppShellHandle {
  showBoard(board: RenderBoard): void;
}

export function mountAppShell(containers: AppShellContainers): AppShellHandle {
  let scene: CubeScene | null = null;
  let current: RenderBoard | null = null;

  const onBoardChange = (board: RenderBoard) => {
    current = board;
    hud.update(board);
  };

  const showBoard = (board: RenderBoard) => {
    menu.hide();
    containers.scene.hidden = false;
    containers.status.hidden = false;

    if (scene) {
      scene.setBoard(board);
    } else {
      scene = createCubeScene(containers.scene, board, onBoardChange);
    }
    scene.start();
  };

  const startMatch = async (option: DifficultyOption) => {
    showBoard(await NewGame(option.value));
  };

  const restart = () => {
    if (!current) {
      return;
    }
    void NewGame(current.difficulty)
      .then(showBoard)
      .catch((err: unknown) => {
        console.error('Não foi possível reiniciar a partida:', err);
      });
  };

  const returnToMenu = () => {
    containers.scene.hidden = true;
    containers.status.hidden = true;
    menu.show();
  };

  const hud: Hud = mountHud(containers.status, { onRestart: restart, onMenu: returnToMenu });
  const menu = mountMenu(containers.menu, startMatch);

  containers.scene.hidden = true;
  containers.status.hidden = true;

  return { showBoard };
}
