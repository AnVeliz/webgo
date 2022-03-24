<script lang="ts">
  import style from "./app.module.css";

  export let timeText: string;

  import { DisableDevTools } from "./services/dev-tools";
  DisableDevTools();
</script>

<main>
  <div class={style.sign}>
    <span class={style.fastFlicker}>{timeText}</span>
  </div>
</main>
