import {
  AmbientLight,
  DirectionalLight,
  PerspectiveCamera,
  Scene,
  WebGLRenderer,
} from 'three';
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls.js';

import { createCellInput, type CellInput } from './cell-input';
import type { RenderBoard } from './cell-materials';
import { buildCube, type CubeMeshes } from './cube-geometry';
import { createHoverHighlight } from './hover-highlight';

const MIN_ZOOM_DISTANCE = 2.2;
const MAX_ZOOM_DISTANCE = 9;

export interface CubeScene {
  start(): void;
  setBoard(board: RenderBoard): void;
  dispose(): void;
}

export function createCubeScene(
  container: HTMLElement,
  board: RenderBoard,
  onBoardChange?: (board: RenderBoard) => void,
): CubeScene {
  const scene = new Scene();

  const camera = new PerspectiveCamera(50, 1, 0.1, 100);
  camera.position.set(2.6, 2.4, 3.4);

  const renderer = new WebGLRenderer({ antialias: true });
  renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
  container.appendChild(renderer.domElement);

  scene.add(new AmbientLight(0xffffff, 1.5));

  const keyLight = new DirectionalLight(0xffffff, 1.6);
  keyLight.position.set(4, 6, 5);
  scene.add(keyLight);

  const rimLight = new DirectionalLight(0xffffff, 0.7);
  rimLight.position.set(-5, -3, -4);
  scene.add(rimLight);

  const controls = new OrbitControls(camera, renderer.domElement);
  controls.enablePan = false;
  controls.enableDamping = true;
  controls.dampingFactor = 0.08;
  controls.rotateSpeed = 0.75;
  controls.zoomSpeed = 0.8;
  controls.minDistance = MIN_ZOOM_DISTANCE;
  controls.maxDistance = MAX_ZOOM_DISTANCE;
  controls.target.set(0, 0, 0);
  controls.update();

  const hover = createHoverHighlight(renderer.domElement, camera);

  let cube: CubeMeshes | null = null;

  const mountBoard = (next: RenderBoard) => {
    if (cube) {
      hover.setTargets([]);
      scene.remove(cube.group);
      cube.dispose();
    }
    cube = buildCube(next);
    scene.add(cube.group);
    hover.setTargets(cube.raycastTargets);
    onBoardChange?.(next);
  };

  mountBoard(board);

  const input: CellInput = createCellInput(renderer.domElement, hover, mountBoard);

  const resize = () => {
    const width = container.clientWidth || 1;
    const height = container.clientHeight || 1;
    renderer.setSize(width, height, false);
    camera.aspect = width / height;
    camera.updateProjectionMatrix();
  };
  resize();
  window.addEventListener('resize', resize);

  let frameId = 0;

  const renderLoop = () => {
    frameId = requestAnimationFrame(renderLoop);
    controls.update();
    hover.update();
    renderer.render(scene, camera);
  };

  return {
    start() {
      if (frameId === 0) {
        resize();
        renderLoop();
      }
    },
    setBoard(next) {
      mountBoard(next);
    },
    dispose() {
      cancelAnimationFrame(frameId);
      frameId = 0;
      window.removeEventListener('resize', resize);
      input.dispose();
      hover.dispose();
      controls.dispose();
      if (cube) {
        scene.remove(cube.group);
        cube.dispose();
        cube = null;
      }
      renderer.dispose();
      renderer.domElement.remove();
    },
  };
}
