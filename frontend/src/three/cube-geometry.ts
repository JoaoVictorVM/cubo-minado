import {
  BufferGeometry,
  Group,
  Matrix4,
  Mesh,
  MeshStandardMaterial,
  PlaneGeometry,
  Vector3,
} from 'three';

import { GRID_COLOR, materialFor, type RenderBoard, type RenderCell } from './cell-materials';

export const CUBE_HALF_EXTENT = 1;

export interface CellIdentity {
  face: number;
  row: number;
  col: number;
}

interface FaceLayout {
  face: number;
  center: Vector3;
  rowDir: Vector3;
  colDir: Vector3;
}

const FACE_LAYOUT: FaceLayout[] = [
  {
    face: 0,
    center: new Vector3(0, CUBE_HALF_EXTENT, 0),
    rowDir: new Vector3(0, 0, 1),
    colDir: new Vector3(1, 0, 0),
  },
  {
    face: 1,
    center: new Vector3(0, -CUBE_HALF_EXTENT, 0),
    rowDir: new Vector3(0, 0, -1),
    colDir: new Vector3(1, 0, 0),
  },
  {
    face: 2,
    center: new Vector3(0, 0, CUBE_HALF_EXTENT),
    rowDir: new Vector3(0, -1, 0),
    colDir: new Vector3(1, 0, 0),
  },
  {
    face: 3,
    center: new Vector3(0, 0, -CUBE_HALF_EXTENT),
    rowDir: new Vector3(0, -1, 0),
    colDir: new Vector3(-1, 0, 0),
  },
  {
    face: 4,
    center: new Vector3(-CUBE_HALF_EXTENT, 0, 0),
    rowDir: new Vector3(0, -1, 0),
    colDir: new Vector3(0, 0, 1),
  },
  {
    face: 5,
    center: new Vector3(CUBE_HALF_EXTENT, 0, 0),
    rowDir: new Vector3(0, -1, 0),
    colDir: new Vector3(0, 0, -1),
  },
];

const CELL_GAP_RATIO = 0.06;
const CELL_SURFACE_OFFSET = 0.004;

function orientationMatrix(layout: FaceLayout): Matrix4 {
  const normal = new Vector3().crossVectors(layout.rowDir, layout.colDir);
  const up = layout.rowDir.clone().negate();
  return new Matrix4().makeBasis(layout.colDir.clone(), up, normal);
}

export interface CubeMeshes {
  group: Group;
  cellMeshes: Mesh[];
  raycastTargets: Mesh[];
  dispose(): void;
}

export function buildCube(board: RenderBoard): CubeMeshes {
  const size = board.boardSize;
  const cellSize = (CUBE_HALF_EXTENT * 2) / size;
  const meshSize = cellSize * (1 - CELL_GAP_RATIO);

  const group = new Group();
  const cellMeshes: Mesh[] = [];
  const backdrops: Mesh[] = [];
  const geometries: BufferGeometry[] = [];
  const ownedMaterials: MeshStandardMaterial[] = [];

  const cellGeometry = new PlaneGeometry(meshSize, meshSize);
  const backdropGeometry = new PlaneGeometry(CUBE_HALF_EXTENT * 2, CUBE_HALF_EXTENT * 2);
  geometries.push(cellGeometry, backdropGeometry);

  for (const layout of FACE_LAYOUT) {
    const faceData = board.faces[layout.face];
    if (!faceData) {
      continue;
    }

    const faceGroup = new Group();
    const orientation = orientationMatrix(layout);
    const normal = new Vector3().crossVectors(layout.rowDir, layout.colDir);

    const backdropMaterial = new MeshStandardMaterial({
      color: GRID_COLOR,
      roughness: 1,
      metalness: 0,
    });
    ownedMaterials.push(backdropMaterial);

    const backdrop = new Mesh(backdropGeometry, backdropMaterial);
    backdrop.position.copy(layout.center);
    backdrop.setRotationFromMatrix(orientation);
    faceGroup.add(backdrop);
    backdrops.push(backdrop);

    for (let row = 0; row < size; row++) {
      for (let col = 0; col < size; col++) {
        const cell: RenderCell | undefined = faceData.cells[row]?.[col];
        if (!cell) {
          continue;
        }

        const rowOffset = (row + 0.5) * cellSize - CUBE_HALF_EXTENT;
        const colOffset = (col + 0.5) * cellSize - CUBE_HALF_EXTENT;

        const position = layout.center
          .clone()
          .addScaledVector(layout.rowDir, rowOffset)
          .addScaledVector(layout.colDir, colOffset)
          .addScaledVector(normal, CELL_SURFACE_OFFSET);

        const mesh = new Mesh(cellGeometry, materialFor(cell, false));
        mesh.position.copy(position);
        mesh.setRotationFromMatrix(orientation);
        mesh.userData.cell = { face: layout.face, row, col } satisfies CellIdentity;
        mesh.userData.appearance = cell;

        faceGroup.add(mesh);
        cellMeshes.push(mesh);
      }
    }

    group.add(faceGroup);
  }

  return {
    group,
    cellMeshes,
    raycastTargets: [...cellMeshes, ...backdrops],
    dispose() {
      for (const geometry of geometries) {
        geometry.dispose();
      }
      for (const material of ownedMaterials) {
        material.dispose();
      }
      group.clear();
    },
  };
}

export function cellIdentityOf(mesh: Mesh): CellIdentity | undefined {
  return mesh.userData.cell as CellIdentity | undefined;
}

export function setCellHighlighted(mesh: Mesh, highlighted: boolean): void {
  const cell = mesh.userData.appearance as RenderCell | undefined;
  if (cell) {
    mesh.material = materialFor(cell, highlighted);
  }
}
