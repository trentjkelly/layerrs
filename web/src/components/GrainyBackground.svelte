<script></script>
  
<div class="absolute w-full h-full grainy-bg opacity-35 z-0"></div>

<svg class="absolute w-0 h-0">
    <filter id="grainy">
      <feTurbulence
        type="turbulence"
        baseFrequency="0.65"
      ></feTurbulence>
    </filter>
</svg>
  
<style>
    .grainy-bg {
        filter: url(#grainy);
    }
</style>