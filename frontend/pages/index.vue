<script setup lang="ts">
definePageMeta({ layout: 'default' })

useHead({ title: 'Veltiq — аналитика чековых книг' })

const features = [
  {
    icon: 'lucide:upload-cloud',
    title: 'Загружай чеки за период',
    text: 'CSV, Excel или выгрузка из 1С — Veltiq распарсит чековую книгу автоматически.',
  },
  {
    icon: 'lucide:line-chart',
    title: 'Видь сводку и динамику',
    text: 'Сумма продаж, доли наличных и карт, средний чек — всё в одном дашборде.',
  },
  {
    icon: 'lucide:package-search',
    title: 'Сокращай неликвид',
    text: 'Получай рекомендации, какие товары стоит распродать, а какие списать.',
  },
]
</script>

<template>
  <section class="mx-auto max-w-6xl px-4 py-16 sm:py-24">
    <div class="grid items-center gap-12 lg:grid-cols-2">
      <div>
        <span class="inline-flex items-center gap-2 rounded-full bg-brand-50 px-3 py-1 text-xs font-medium text-brand-700">
          <Icon name="lucide:sparkles" class="size-3.5" />
          SaaS для малой розницы
        </span>
        <h1 class="mt-6 text-4xl font-semibold tracking-tight text-surface-900 sm:text-5xl">
          Сокращайте неликвид
          <span class="text-brand-600">через аналитику чеков</span>
        </h1>
        <p class="mt-5 max-w-xl text-lg text-surface-600">
          Veltiq разбирает чековые книги ваших магазинов и подсказывает,
          какие товары давно не продаются и что с ними делать.
        </p>
        <div class="mt-8 flex flex-wrap gap-3">
          <VButton to="/register" size="lg" icon-right="lucide:arrow-right">
            Попробовать бесплатно
          </VButton>
          <VButton to="/login" size="lg" variant="outline">Войти</VButton>
        </div>
      </div>

      <div class="card-base p-6">
        <div class="grid grid-cols-2 gap-4 text-sm">
          <div class="rounded-xl bg-surface-50 p-4">
            <div class="text-xs text-surface-500">Чеков за период</div>
            <div class="mt-1 text-2xl font-semibold text-surface-900">12 480</div>
          </div>
          <div class="rounded-xl bg-surface-50 p-4">
            <div class="text-xs text-surface-500">Оборот</div>
            <div class="mt-1 text-2xl font-semibold text-surface-900">3,2 млн ₽</div>
          </div>
          <div class="rounded-xl bg-surface-50 p-4 col-span-2">
            <div class="text-xs text-surface-500">Кандидаты на распродажу</div>
            <ul class="mt-2 space-y-1.5 text-surface-700">
              <li class="flex items-center justify-between">
                <span>Чай зелёный, 100г</span>
                <VBadge tone="warning" icon="lucide:timer">62 дня без продаж</VBadge>
              </li>
              <li class="flex items-center justify-between">
                <span>Печенье «Утро», 200г</span>
                <VBadge tone="danger" icon="lucide:trending-down">списать</VBadge>
              </li>
              <li class="flex items-center justify-between">
                <span>Сок «Луч», 1л</span>
                <VBadge tone="info">скидка 15%</VBadge>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </div>

    <div class="mt-20 grid gap-6 sm:grid-cols-3">
      <VCard
        v-for="f in features"
        :key="f.title"
        :title="f.title"
        :description="f.text"
      >
        <template #actions>
          <div class="flex size-9 items-center justify-center rounded-xl bg-brand-50 text-brand-600">
            <Icon :name="f.icon" class="size-5" />
          </div>
        </template>
      </VCard>
    </div>
  </section>
</template>
