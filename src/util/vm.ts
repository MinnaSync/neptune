import vm from 'node:vm';

const context = vm.createContext({});

/**
 * Runs code in a sandboxed environment.
 * This should not be considered a "complete" security mechanism.
 * It's only used to prevent access to the main process.
 * @see https://nodejs.org/api/vm.html#vm-executing-javascript
 */
export function evalSync(code: string) {
    try {
        return vm.runInContext(code, context);
    } catch (e) {
        return null;
    }
}