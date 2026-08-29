import './style.css';

import { mountAppShell } from './app-shell';
import { mountDebugPanel } from './debug-panel';

const sceneContainer = document.getElementById('scene');
const menuContainer = document.getElementById('menu');
const statusContainer = document.getElementById('game-status');
const debugContainer = document.getElementById('debug-panel');

if (!sceneContainer || !menuContainer || !statusContainer || !debugContainer) {
  throw new Error('Elementos de montagem não encontrados no index.html');
}

const shell = mountAppShell({
  menu: menuContainer,
  scene: sceneContainer,
  status: statusContainer,
});

mountDebugPanel(debugContainer, shell);
