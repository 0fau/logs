export function horizontalWheel(container: HTMLElement) {
    /** Max `scrollLeft` value */
    let scrollWidth: number;

    /** Desired scroll distance per animation frame */
    let getScrollStep = () => scrollWidth / 50; /* ADJUST TO YOUR WISH */

    /** Target value for `scrollLeft` */
    let targetLeft: number;

    function scrollLeft() {
        let beforeLeft = container.scrollLeft;
        let wantDx = getScrollStep();
        let diff = targetLeft - container.scrollLeft;
        let dX = wantDx >= Math.abs(diff) ? diff : Math.sign(diff) * wantDx;

        // Performing horizontal scroll
        container.scrollBy(dX, 0);

        // Break if smaller `diff` instead of `wantDx` was used
        if (dX === diff) return;

        // Break if can't scroll anymore or target reached
        if (beforeLeft === container.scrollLeft || container.scrollLeft === targetLeft) return;

        requestAnimationFrame(scrollLeft);
    }

    container.addEventListener("wheel", (e) => {
        e.preventDefault();

        scrollWidth = container.scrollWidth - container.clientWidth;
        targetLeft = Math.min(scrollWidth, Math.max(0, container.scrollLeft + e.deltaY));

        requestAnimationFrame(scrollLeft);
    });
}
