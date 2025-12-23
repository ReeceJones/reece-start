<script lang="ts">
	import { onMount, tick } from 'svelte';
	import type { CropperCanvas, CropperImage, CropperSelection } from 'cropperjs';
	import { Button } from '../ui/button';

	let {
		imageFile,
		onCancel,
		onSave
	}: { imageFile: File; onCancel: () => void; onSave: (file: File) => void } = $props();

	let canvasRef: CropperCanvas | null = null;
	let imageRef: CropperImage | null = null;
	let cropperSelection: CropperSelection | null = null;
	let imageLoaded = $state(false);

	const imageUrl = $derived(URL.createObjectURL(imageFile));

	async function onTransformChange() {
		// The transform event is fired before the transform is applied so we need to wait a tick
		// to ensure the selection constraint is applied to the correct transform state.
		await tick();
		if (cropperSelection) {
			constrainSelection({
				x: cropperSelection.x,
				y: cropperSelection.y,
				width: cropperSelection.width,
				height: cropperSelection.height
			});
		}
	}

	function onSelectionChange(event: CustomEvent) {
		// Only apply constraints if the image is loaded
		if (!imageLoaded) {
			return;
		}

		// Prevent the default behavior of the selection change event
		// to avoid the cropper from applying the selection immediately.
		// We want to apply our own constraints instead.
		event.preventDefault();

		constrainSelection(event.detail);
	}

	function constrainSelection(selection: { x: number; y: number; width: number; height: number }) {
		if (!canvasRef || !imageRef || !cropperSelection) {
			return;
		}

		const cropperCanvasRect = canvasRef.getBoundingClientRect();
		const cropperImageRect = imageRef.getBoundingClientRect();

		// Don't apply constraints if the image hasn't been sized yet
		if (cropperImageRect.width === 0 || cropperImageRect.height === 0) {
			return;
		}

		// Define minimum selection size (10% of image dimensions)
		const minWidth = Math.max(20, cropperImageRect.width * 0.1);
		const minHeight = Math.max(20, cropperImageRect.height * 0.1);

		// Ensure selection has minimum dimensions
		let width = Math.max(minWidth, Math.min(selection.width, cropperImageRect.width));
		let height = Math.max(minHeight, Math.min(selection.height, cropperImageRect.height));

		// For square aspect ratio, use the smaller dimension
		if (cropperSelection.getAttribute('aspect-ratio') === '1') {
			const size = Math.min(width, height);
			width = size;
			height = size;
		}

		const minX = cropperImageRect.x - cropperCanvasRect.x;
		const minY = cropperImageRect.y - cropperCanvasRect.y;
		const maxX = minX + cropperImageRect.width - width;
		const maxY = minY + cropperImageRect.height - height;

		const x = Math.max(minX, Math.min(selection.x, maxX));
		const y = Math.max(minY, Math.min(selection.y, maxY));

		// Only update if the values have actually changed to avoid unnecessary renders
		if (
			cropperSelection.x !== x ||
			cropperSelection.y !== y ||
			cropperSelection.width !== width ||
			cropperSelection.height !== height
		) {
			cropperSelection.x = x;
			cropperSelection.y = y;
			cropperSelection.width = width;
			cropperSelection.height = height;
			cropperSelection.$render();
		}
	}

	function maximizeSelectionForSquareImage() {
		if (!canvasRef || !imageRef || !cropperSelection) {
			return;
		}

		const cropperCanvasRect = canvasRef.getBoundingClientRect();
		const cropperImageRect = imageRef.getBoundingClientRect();

		// Don't apply if the image hasn't been sized yet
		if (cropperImageRect.width === 0 || cropperImageRect.height === 0) {
			return;
		}

		// Check if the image is square (or very close to square, within 1% tolerance)
		const aspectRatio = cropperImageRect.width / cropperImageRect.height;
		const isSquare = Math.abs(aspectRatio - 1) < 0.01;

		if (isSquare) {
			console.log('Image is square, maximizing selection to cover entire image');

			// Calculate position relative to canvas
			const x = cropperImageRect.x - cropperCanvasRect.x;
			const y = cropperImageRect.y - cropperCanvasRect.y;

			// Set selection to cover the entire image
			cropperSelection.x = x;
			cropperSelection.y = y;
			cropperSelection.width = cropperImageRect.width;
			cropperSelection.height = cropperImageRect.height;
			cropperSelection.$render();
		}
	}

	onMount(async () => {
		await import('cropperjs');
	});

	// Set up the ready callback when imageRef becomes available
	$effect(() => {
		if (imageRef) {
			console.log('imageRef is available:', imageRef);

			// Function to mark image as ready
			const markImageReady = () => {
				if (!imageLoaded) {
					console.log('Image marked as ready');
					imageLoaded = true;

					// center the viewport on the image
					imageRef?.$center('contain');

					// Check if image is square and maximize selection if so
					maximizeSelectionForSquareImage();
				}
			};

			// Check if image is already ready (race condition fix)
			try {
				// @ts-expect-error CropperJS may have different ready state properties
				if (imageRef.ready || imageRef.$ready?.length === 0) {
					console.log('Image was already ready');
					markImageReady();
					return;
				}
			} catch (e: unknown) {
				console.error('Error checking ready state:', e);
			}

			// Try different approaches to detect when image becomes ready
			try {
				// Method 1: Try $ready
				if (typeof imageRef.$ready === 'function') {
					console.log('Using $ready method');
					imageRef.$ready(markImageReady);
				} else {
					console.log('$ready method not available, trying addEventListener');
					// Method 2: Try ready event
					imageRef.addEventListener('ready', markImageReady);
				}

				// Fallback: Set a timeout as backup
				setTimeout(() => {
					if (!imageLoaded) {
						console.log('Image ready via timeout fallback');
						markImageReady();
					}
				}, 500); // Increased timeout to be more conservative
			} catch (error) {
				console.error('Error setting up ready callback:', error);
				// Immediate fallback if all methods fail
				setTimeout(markImageReady, 100);
			}
		}
	});
</script>

<div class="flex flex-col gap-2">
	<div class="aspect-square w-full">
		<cropper-canvas bind:this={canvasRef} background class="aspect-square">
			<cropper-image
				bind:this={imageRef}
				src={imageUrl}
				alt="uploaded picture"
				scalable
				translatable
				ontransform={onTransformChange}
			></cropper-image>
			<cropper-shade hidden></cropper-shade>
			<cropper-handle action="select" plain></cropper-handle>
			<cropper-selection
				bind:this={cropperSelection}
				initial-coverage="0.8"
				movable
				resizable
				aspect-ratio="1"
				precise
				id="logo-cropper-selection"
				onchange={onSelectionChange}
			>
				<cropper-grid role="grid" covered></cropper-grid>
				<cropper-crosshair centered></cropper-crosshair>
				<cropper-handle action="move" theme-color="rgba(255, 255, 255, 0.35)"></cropper-handle>
				<cropper-handle action="n-resize"></cropper-handle>
				<cropper-handle action="e-resize"></cropper-handle>
				<cropper-handle action="s-resize"></cropper-handle>
				<cropper-handle action="w-resize"></cropper-handle>
				<cropper-handle action="ne-resize"></cropper-handle>
				<cropper-handle action="nw-resize"></cropper-handle>
				<cropper-handle action="se-resize"></cropper-handle>
				<cropper-handle action="sw-resize"></cropper-handle>
			</cropper-selection>
		</cropper-canvas>
	</div>

	<div class="flex justify-end gap-2">
		<Button variant="outline" type="button" onclick={() => onCancel()}>Cancel</Button>
		<Button
			type="button"
			onclick={async () => {
				// get the cropped image
				const cropperSelection = document.getElementById('logo-cropper-selection');
				if (!cropperSelection) return;
				console.log(cropperSelection);
				// @ts-expect-error cropper-selection is not typed
				const selectionCanvas = await cropperSelection.$toCanvas();
				console.log(selectionCanvas);

				// Convert canvas to blob first, then to file
				selectionCanvas.toBlob((blob: Blob | null) => {
					if (blob) {
						console.log('Blob created:', blob);
						console.log('Blob size:', blob.size);
						console.log('Blob type:', blob.type);

						// Create a proper file from the blob
						const file = new File([blob], imageFile.name, {
							type: 'image/png'
						});
						console.log('File created:', file);
						console.log('File size:', file.size);
						console.log('File type:', file.type);

						onSave(file);
					} else {
						console.error('Failed to create blob from canvas');
					}
				}, 'image/png');
			}}
		>
			Save
		</Button>
	</div>

	<!-- Render the canvas element returned by the cropper -->
	<!-- <canvas bind:this={canvas}></canvas> -->
</div>
