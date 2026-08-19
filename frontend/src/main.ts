import './style.css';

import { mountDebugPanel } from './debug-panel';
import { createPlaceholderScene } from './three/scene';

const sceneContainer = document.getElementById('scene');
const debugContainer = document.getElementById('debug-panel');

if (!sceneContainer || !debugContainer) {
  throw new Error('Elementos de montagem não encontrados no index.html');
}

const scene = createPlaceholderScene(sceneContainer);
scene.start();

mountDebugPanel(debugContainer);
