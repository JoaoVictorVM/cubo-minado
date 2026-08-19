import {
  BoxGeometry,
  Mesh,
  MeshStandardMaterial,
  PerspectiveCamera,
  DirectionalLight,
  AmbientLight,
  Scene,
  WebGLRenderer,
} from 'three';

export interface PlaceholderScene {
  start(): void;
  dispose(): void;
}

export function createPlaceholderScene(container: HTMLElement): PlaceholderScene {
  const scene = new Scene();

  const camera = new PerspectiveCamera(60, 1, 0.1, 100);
  camera.position.set(2.5, 2.5, 3.5);
  camera.lookAt(0, 0, 0);

  const renderer = new WebGLRenderer({ antialias: true });
  renderer.setPixelRatio(window.devicePixelRatio);
  container.appendChild(renderer.domElement);

  scene.add(new AmbientLight(0xffffff, 0.6));
  const keyLight = new DirectionalLight(0xffffff, 1.2);
  keyLight.position.set(3, 4, 5);
  scene.add(keyLight);

  const cube = new Mesh(
    new BoxGeometry(2, 2, 2),
    new MeshStandardMaterial({ color: 0x4d9de0, roughness: 0.4, metalness: 0.1 }),
  );
  scene.add(cube);

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
    cube.rotation.x += 0.005;
    cube.rotation.y += 0.008;
    renderer.render(scene, camera);
  };

  return {
    start() {
      if (frameId === 0) {
        renderLoop();
      }
    },
    dispose() {
      cancelAnimationFrame(frameId);
      frameId = 0;
      window.removeEventListener('resize', resize);
      cube.geometry.dispose();
      (cube.material as MeshStandardMaterial).dispose();
      renderer.dispose();
      renderer.domElement.remove();
    },
  };
}
