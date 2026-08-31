import { ChordCell, FlagCell, OpenCell } from '../../wailsjs/go/main/App';
import type { RenderBoard } from './cell-materials';
import type { CellIdentity } from './cube-geometry';
import type { HoverHighlight } from './hover-highlight';

const LEFT_BUTTON = 0;
const RIGHT_BUTTON = 2;
const CLICK_DRAG_THRESHOLD_PX = 5;

type CellAction = (face: number, row: number, col: number) => Promise<RenderBoard>;

export interface CellInput {
  dispose(): void;
}

export function createCellInput(
  canvas: HTMLCanvasElement,
  hover: HoverHighlight,
  onBoardUpdate: (board: RenderBoard) => void,
): CellInput {
  const heldButtons = new Set<number>();
  let pointerDownAt: { x: number; y: number } | null = null;
  let chordFiredForGesture = false;

  const canvasPoint = (event: PointerEvent) => {
    const rect = canvas.getBoundingClientRect();
    return { x: event.clientX - rect.left, y: event.clientY - rect.top };
  };

  const movedPastThreshold = (event: PointerEvent) => {
    if (!pointerDownAt) {
      return true;
    }
    const now = canvasPoint(event);
    return Math.hypot(now.x - pointerDownAt.x, now.y - pointerDownAt.y) > CLICK_DRAG_THRESHOLD_PX;
  };

  const resetGesture = () => {
    heldButtons.clear();
    pointerDownAt = null;
    chordFiredForGesture = false;
  };

  const dispatch = (action: CellAction, target: CellIdentity) => {
    action(target.face, target.row, target.col)
      .then(onBoardUpdate)
      .catch((err: unknown) => {
        const message = err instanceof Error ? err.message : String(err);
        if (message === 'ERR_NOT_IMPLEMENTED') {
          return;
        }
        console.error(`Ação de célula falhou em (${target.face},${target.row},${target.col}):`, message);
      });
  };

  const onPointerDown = (event: PointerEvent) => {
    if (event.button !== LEFT_BUTTON && event.button !== RIGHT_BUTTON) {
      return;
    }

    if (heldButtons.size === 0) {
      pointerDownAt = canvasPoint(event);
      chordFiredForGesture = false;
    }
    heldButtons.add(event.button);

    if (
      !chordFiredForGesture &&
      heldButtons.has(LEFT_BUTTON) &&
      heldButtons.has(RIGHT_BUTTON)
    ) {
      chordFiredForGesture = true;
      const target = hover.getHighlighted();
      if (target) {
        dispatch(ChordCell, target);
      }
    }
  };

  const onPointerUp = (event: PointerEvent) => {
    if (!heldButtons.delete(event.button)) {
      return;
    }

    if (!chordFiredForGesture && !movedPastThreshold(event)) {
      const target = hover.getHighlighted();
      if (target) {
        dispatch(event.button === RIGHT_BUTTON ? FlagCell : OpenCell, target);
      }
    }

    if (heldButtons.size === 0) {
      resetGesture();
    }
  };

  canvas.addEventListener('pointerdown', onPointerDown);
  canvas.addEventListener('pointerup', onPointerUp);
  canvas.addEventListener('pointerleave', resetGesture);

  return {
    dispose() {
      resetGesture();
      canvas.removeEventListener('pointerdown', onPointerDown);
      canvas.removeEventListener('pointerup', onPointerUp);
      canvas.removeEventListener('pointerleave', resetGesture);
    },
  };
}
