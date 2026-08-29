import type { RenderBoard, RenderCell } from './cell-materials';

const PREVIEW_SIZE = 5;
const FRONT_FACE = 2;

function closedCell(): RenderCell {
  return { state: 'closed', adjacentMines: 0 };
}

function numberedCell(count: number): RenderCell {
  return { state: 'open', adjacentMines: count };
}

function emptyFace(size: number): RenderCell[][] {
  return Array.from({ length: size }, () =>
    Array.from({ length: size }, () => closedCell()),
  );
}

export function buildCellStatePreviewBoard(): RenderBoard {
  const faces = Array.from({ length: 6 }, () => ({ cells: emptyFace(PREVIEW_SIZE) }));

  const front = faces[FRONT_FACE].cells;

  front[0] = [
    numberedCell(1),
    numberedCell(2),
    numberedCell(3),
    numberedCell(4),
    numberedCell(5),
  ];

  front[1] = [
    numberedCell(6),
    numberedCell(7),
    numberedCell(8),
    { state: 'flagged', adjacentMines: 0 },
    { state: 'open', adjacentMines: 0, mine: true, triggered: true },
  ];

  front[2] = [
    { state: 'open', adjacentMines: 0, mine: true },
    { state: 'open', adjacentMines: 0 },
    closedCell(),
    { state: 'flagged', adjacentMines: 0 },
    closedCell(),
  ];

  return { boardSize: PREVIEW_SIZE, faces };
}
