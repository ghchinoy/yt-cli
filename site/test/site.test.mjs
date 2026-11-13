import { test } from 'node:test';
import assert from 'node:assert/strict';
import { readFileSync, existsSync } from 'node:fs';
import { join, dirname } from 'node:path';
import { fileURLToPath } from 'node:url';

const __dirname = dirname(fileURLToPath(import.meta.url));
const distIndex = join(__dirname, '..', 'dist', 'index.html');

test('built index.html contains core yt-cli sections and zero external font/tracker requests', () => {
  assert.ok(existsSync(distIndex), 'dist/index.html must exist (run npm run build first)');
  const html = readFileSync(distIndex, 'utf8');

  for (const expected of [
    'YouTube CLI',
    'Rich for humans.',
    'Deterministic for agents.',
    'yt-cli channel info',
    'yt-cli channel videos',
    'yt-cli video channel-info',
    'yt-cli video upload',
    'go install github.com/ghchinoy/yt-cli@latest',
    'gemini skills install https://github.com/ghchinoy/yt-cli.git --path skills/yt-cli',
  ]) {
    assert.ok(html.includes(expected), `expected dist/index.html to include "${expected}"`);
  }

  assert.ok(
    !html.includes('fonts.googleapis.com'),
    'dist/index.html must not request fonts from fonts.googleapis.com'
  );
});
