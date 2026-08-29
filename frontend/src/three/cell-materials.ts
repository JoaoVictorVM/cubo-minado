import {
  CanvasTexture,
  LinearFilter,
  MeshStandardMaterial,
  SRGBColorSpace,
  type Texture,
} from 'three';

export const CLOSED_COLOR = '#b9b9b9';
export const OPENED_COLOR = '#dedad2';
export const TRIGGERED_MINE_COLOR = '#d1495b';
export const GRID_COLOR = '#4a4a4a';

const NUMBER_COLORS: Record<number, string> = {
  1: '#0000ff',
  2: '#008000',
  3: '#ff0000',
  4: '#000080',
  5: '#800000',
  6: '#008080',
  7: '#000000',
  8: '#808080',
};

const TEXTURE_SIZE = 128;

const textureCache = new Map<string, Texture>();
const materialCache = new Map<string, MeshStandardMaterial>();

export interface RenderCell {
  state: string;
  adjacentMines: number;
  mine?: boolean;
  triggered?: boolean;
}

export interface RenderFace {
  cells: RenderCell[][];
}

export interface RenderBoard {
  boardSize: number;
  faces: RenderFace[];
}

function drawTexture(key: string, paint: (ctx: CanvasRenderingContext2D, size: number) => void): Texture {
  const cached = textureCache.get(key);
  if (cached) {
    return cached;
  }

  const canvas = document.createElement('canvas');
  canvas.width = TEXTURE_SIZE;
  canvas.height = TEXTURE_SIZE;
  const ctx = canvas.getContext('2d');
  if (!ctx) {
    throw new Error('Canvas 2D indisponível para gerar as texturas das células');
  }

  paint(ctx, TEXTURE_SIZE);

  const texture = new CanvasTexture(canvas);
  texture.colorSpace = SRGBColorSpace;
  texture.anisotropy = 4;
  texture.minFilter = LinearFilter;
  textureCache.set(key, texture);
  return texture;
}

function fill(ctx: CanvasRenderingContext2D, size: number, color: string): void {
  ctx.fillStyle = color;
  ctx.fillRect(0, 0, size, size);
}

function numberTexture(count: number): Texture {
  return drawTexture(`number-${count}`, (ctx, size) => {
    fill(ctx, size, OPENED_COLOR);
    ctx.fillStyle = NUMBER_COLORS[count];
    ctx.font = `bold ${size * 0.72}px "Segoe UI", system-ui, sans-serif`;
    ctx.textAlign = 'center';
    ctx.textBaseline = 'middle';
    ctx.fillText(String(count), size / 2, size * 0.54);
  });
}

function flagTexture(): Texture {
  return drawTexture('flag', (ctx, size) => {
    fill(ctx, size, CLOSED_COLOR);

    ctx.fillStyle = '#1c1c1c';
    ctx.fillRect(size * 0.55, size * 0.22, size * 0.06, size * 0.5);
    ctx.fillRect(size * 0.34, size * 0.72, size * 0.42, size * 0.08);
    ctx.fillRect(size * 0.42, size * 0.66, size * 0.26, size * 0.07);

    ctx.fillStyle = '#d40000';
    ctx.beginPath();
    ctx.moveTo(size * 0.55, size * 0.2);
    ctx.lineTo(size * 0.55, size * 0.5);
    ctx.lineTo(size * 0.24, size * 0.35);
    ctx.closePath();
    ctx.fill();
  });
}

function mineTexture(background: string): Texture {
  return drawTexture(`mine-${background}`, (ctx, size) => {
    fill(ctx, size, background);

    const center = size / 2;
    const radius = size * 0.24;

    ctx.strokeStyle = '#101010';
    ctx.lineWidth = size * 0.07;
    for (let i = 0; i < 4; i++) {
      const angle = (Math.PI / 4) * i;
      const reach = radius * 1.72;
      ctx.beginPath();
      ctx.moveTo(center - Math.cos(angle) * reach, center - Math.sin(angle) * reach);
      ctx.lineTo(center + Math.cos(angle) * reach, center + Math.sin(angle) * reach);
      ctx.stroke();
    }

    ctx.fillStyle = '#101010';
    ctx.beginPath();
    ctx.arc(center, center, radius, 0, Math.PI * 2);
    ctx.fill();

    ctx.fillStyle = '#f2f2f2';
    ctx.beginPath();
    ctx.arc(center - radius * 0.34, center - radius * 0.34, radius * 0.24, 0, Math.PI * 2);
    ctx.fill();
  });
}

export function appearanceKey(cell: RenderCell): string {
  if (cell.state === 'flagged') {
    return 'flagged';
  }
  if (cell.state === 'open') {
    if (cell.mine) {
      return cell.triggered ? 'mine-triggered' : 'mine-revealed';
    }
    if (cell.adjacentMines > 0) {
      return `number-${cell.adjacentMines}`;
    }
    return 'open-empty';
  }
  return 'closed';
}

function buildMaterial(key: string, highlighted: boolean): MeshStandardMaterial {
  const material = new MeshStandardMaterial({ roughness: 0.85, metalness: 0.02 });

  if (key === 'closed') {
    material.color.set(CLOSED_COLOR);
  } else if (key === 'open-empty') {
    material.color.set(OPENED_COLOR);
  } else if (key === 'flagged') {
    material.map = flagTexture();
  } else if (key === 'mine-triggered') {
    material.map = mineTexture(TRIGGERED_MINE_COLOR);
  } else if (key === 'mine-revealed') {
    material.map = mineTexture(OPENED_COLOR);
  } else {
    material.map = numberTexture(Number(key.slice('number-'.length)));
  }

  if (highlighted) {
    material.emissive.set('#4d9de0');
    material.emissiveIntensity = 0.55;
  }

  return material;
}

export function materialFor(cell: RenderCell, highlighted: boolean): MeshStandardMaterial {
  const key = appearanceKey(cell);
  const cacheKey = highlighted ? `${key}:hover` : key;

  const cached = materialCache.get(cacheKey);
  if (cached) {
    return cached;
  }

  const material = buildMaterial(key, highlighted);
  materialCache.set(cacheKey, material);
  return material;
}
