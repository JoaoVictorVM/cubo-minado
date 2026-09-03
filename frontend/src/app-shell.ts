import { GetBestTimes, NewGame, SubmitTime } from '../wailsjs/go/main/App';
import { mountEndScreen, type EndScreen } from './end-screen';
import { elapsedSeconds, mountHud, type Hud } from './hud';
import { mountMenu, type DifficultyOption } from './menu';
import type { RenderBoard } from './three/cell-materials';
import { createCubeScene, type CubeScene } from './three/scene';

export interface AppShellContainers {
  menu: HTMLElement;
  scene: HTMLElement;
  status: HTMLElement;
  endScreen: HTMLElement;
}

export interface AppShellHandle {
  showBoard(board: RenderBoard): void;
}

function bestTimeFor(times: { easy?: number; medium?: number; hard?: number }, difficulty: string) {
  const value = times[difficulty as 'easy' | 'medium' | 'hard'];
  return value === undefined || value === null ? null : value;
}

export function mountAppShell(containers: AppShellContainers): AppShellHandle {
  let scene: CubeScene | null = null;
  let current: RenderBoard | null = null;
  let endHandled = false;

  const closeEndScreen = () => {
    endHandled = false;
    endScreen.hide();
    hud.setActionsHidden(false);
  };

  const revealVictory = async (board: RenderBoard) => {
    const finalSeconds = elapsedSeconds(board);

    let previousBest: number | null = null;
    try {
      previousBest = bestTimeFor(await GetBestTimes(), board.difficulty);
    } catch (err: unknown) {
      console.error('Não foi possível ler o recorde anterior:', err);
    }

    try {
      await SubmitTime(board.difficulty, finalSeconds);
    } catch (err: unknown) {
      console.error('Não foi possível salvar o tempo da partida:', err);
    }

    endScreen.show({
      result: 'victory',
      finalSeconds,
      previousBest,
      isNewRecord: previousBest === null || finalSeconds < previousBest,
    });
  };

  const handleResult = (board: RenderBoard) => {
    if (endHandled || !board.result) {
      return;
    }
    endHandled = true;
    hud.setActionsHidden(true);

    if (board.result === 'victory') {
      void revealVictory(board);
    } else {
      endScreen.show({ result: 'defeat' });
    }
  };

  const onBoardChange = (board: RenderBoard) => {
    current = board;
    hud.update(board);
    handleResult(board);
  };

  const showBoard = (board: RenderBoard) => {
    closeEndScreen();
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
    closeEndScreen();
    containers.scene.hidden = true;
    containers.status.hidden = true;
    menu.show();
  };

  const hud: Hud = mountHud(containers.status, { onRestart: restart, onMenu: returnToMenu });
  const endScreen: EndScreen = mountEndScreen(containers.endScreen, {
    onPlayAgain: restart,
    onMainMenu: returnToMenu,
  });
  const menu = mountMenu(containers.menu, startMatch);

  containers.scene.hidden = true;
  containers.status.hidden = true;

  return { showBoard };
}
