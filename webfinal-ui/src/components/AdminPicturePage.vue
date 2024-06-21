<script setup>
import {ref, onMounted} from 'vue';

const pictures = ref([]);

async function fetchPictures() {
  const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/api/v1/pic/list/all`, {
    headers: {
      'Authorization': 'Bearer ' + localStorage.getItem('token')
    }
  });
  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }
  const data = await response.json();
  if (data && data.data) {
    pictures.value = data.data;
  }
}

async function deletePicture(uuid) {
  const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/api/v1/pic/delete/${uuid}`, {
    headers: {
      'Authorization': 'Bearer ' + localStorage.getItem('token')
    }
  });
  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }
  // Remove the picture from the list
  pictures.value = pictures.value.filter(picture => picture.UUID !== uuid);
}

async function previewPicture(uuid) {
  window.open(`${import.meta.env.VITE_BACKEND_URL}/img/` + uuid, '_blank');
}

async function downloadPicture(uuid) {
  const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/api/v1/pic/download/${uuid}`, {
    headers: {
      'Authorization': 'Bearer ' + localStorage.getItem('token')
    }
  });

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }

  // Get the filename from the Content-Disposition header
  const contentDisposition = response.headers.get('Content-Disposition');
  let filename = 'download';
  if (contentDisposition) {
    const filenameRegex = /filename[^;=\n]*=((['"]).*?\2|[^;\n]*)/;
    let matches = filenameRegex.exec(contentDisposition);
    if (matches != null && matches[1]) {
      filename = matches[1].replace(/['"]/g, '');
    }
  }

  // Get the blob data
  const blob = await response.blob();

  // Create a link element
  const url = window.URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.setAttribute('download', filename);
  document.body.appendChild(link);
  link.click();

  // Clean up
  document.body.removeChild(link);
  window.URL.revokeObjectURL(url);
}

onMounted(fetchPictures);
</script>

<template>
  <mdui-list class="picture-list">
    <mdui-list-item v-for="(picture, index) in pictures" :key="index"
                    :headline="picture.Name" :description="'Owner ID: '+picture.UserID">
      <mdui-button variant="elevated" slot="end-icon" @click="previewPicture(picture.UUID)">Preview
      </mdui-button>
      <mdui-button variant="elevated" slot="end-icon" @click="downloadPicture(picture.UUID)">Download
      </mdui-button>
      <mdui-button variant="elevated" slot="end-icon" @click="deletePicture(picture.UUID)">Delete
      </mdui-button>
    </mdui-list-item>
  </mdui-list>
</template>

<style scoped>


.picture-list {
  top: 0;
}
</style>