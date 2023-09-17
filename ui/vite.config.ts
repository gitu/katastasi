import {fileURLToPath, URL} from 'node:url'

import {ConfigEnv, defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import {execSync} from "node:child_process";

export default ({mode}: ConfigEnv) => {

    const commitDate = execSync('git log -1 --format=%cI').toString().trimEnd();
    const commitVersion = execSync('git describe --tags  --match "v*" --always HEAD').toString().trimEnd();
    const commitHash = execSync('git rev-parse HEAD').toString().trimEnd();
    const lastCommitMessage = execSync('git show -s --format=%s').toString().trimEnd();

    process.env.VITE_GIT_COMMIT_DATE = commitDate;
    process.env.VITE_GIT_COMMIT_VERSION = commitVersion;
    process.env.VITE_GIT_COMMIT_HASH = commitHash;
    process.env.VITE_GIT_LAST_COMMIT_MESSAGE = lastCommitMessage;

    return defineConfig({
        plugins: [vue(), vueJsx()],
        resolve: {
            alias: {
                '@': fileURLToPath(new URL('./src', import.meta.url))
            }
        },

        server: {
            proxy: {
                "/api": {
                    target: "http://localhost:1323",
                },
            },
        }
    });
};